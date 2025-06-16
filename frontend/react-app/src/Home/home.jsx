import { useEffect, useReducer, useState } from "react";
import { useLocation } from "react-router";
import { useNavigate } from "react-router";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  // faCircleUser,
  faFolderPlus,
  faPenToSquare,
  faTrashCan,
  faArrowsLeftRight,
  faXmark,
} from "@fortawesome/free-solid-svg-icons";
import { MilkdownProvider } from "@milkdown/react";
import PropTypes from "prop-types";

import HomeEditor from "../components/editor";
import PreHome from "./prehome.jsx";
import Base from "../Index/base.jsx";
import {
  ErrorModal,
  FileFolderCreateModal,
  FileFolderDeleteModal,
} from "../components/Modal.jsx";

const locationPreHomePath = "/prehome";
const locationHomePath = "/home";
const apiHomePath = "/api/home";
const apiPreHomePath = "/api/home/prehome";
const apiCreateNewFilePath = "/api/home/create-new-file";
const apiCreateNewFolderPath = "/api/home/create-new-folder";
const apiDeleteFilePath = "/api/home/delete-file";
const apiDeleteFolderPath = "/api/home/delete-folder";

const Home = () => {
  const [data, setData] = useState({
    user_id: "",
    username: "",
    email: "",
    errMsg: "",
    phase: "",
    projects: [],
    project: {},
    projectFolders: [],
    filesPerFolder: {},
  });
  const [status, setStatus] = useState(0);
  const [toggleSidebar, setToggleSidebar] = useState(true);
  const [width, setWidth] = useState(200);
  const [projectPerFolder, setProjectPerFolder] = useState([]);
  const [openFolder, setOpenFolder] = useState({});
  const [selectFile, setSelectFile] = useState(null);
  const [selectFolder, setSelectFolder] = useState(null);
  const [isOpenModal, setIsOpenModal] = useState(false);
  const [errors, setErrors] = useState("");
  const location = useLocation();

  const navigate = useNavigate();

  const initialState = {
    openFile: {},
    currentFile: null,
  };

  const initialOpenModalState = {
    operationType: null,
    targetType: null,
  };

  function reducer(state, action) {
    switch (action.type) {
      case "close_file": {
        console.log("close file", action.payload);

        const { filename } = action.payload;
        const { [filename]: _, ...newOpenFile } = state.openFile;

        const openFileKeysList = Object.keys(state.openFile);
        const currentOpenFileCount = Object.keys(newOpenFile).length;
        let nextFileContent;

        for (let i = 0; i < openFileKeysList.length; i++) {
          if (currentOpenFileCount === 0) {
            nextFileContent = null;
          }

          if (currentOpenFileCount >= 1 && filename === openFileKeysList[i]) {
            if (i === 0 && currentOpenFileCount >= 1) {
              nextFileContent = newOpenFile[openFileKeysList[i + 1]];
            } else if (i >= 1) {
              nextFileContent = newOpenFile[openFileKeysList[i - 1]];
            }
          }
        }

        return {
          ...state,
          openFile: newOpenFile,
          currentFile: nextFileContent,
        };
      }

      case "open_file": {
        const { id, filename } = action.payload;
        let newFile;
        const foldersFileList = Object.values(data.filesPerFolder);
        foldersFileList.forEach((files) => {
          files.forEach((file) => {
            if (file.file_id === id && file.filename === filename) {
              newFile = file;
            }
          });
        });
        console.log("current file: ", newFile);
        return {
          ...state,
          openFile: {
            ...state.openFile,
            [filename]: newFile,
          },
          currentFile: newFile,
        };
      }

      case "update_file": {
        const { file } = action.payload;
        console.log("update file: ", file);
        return {
          ...state,
          currentFile: file,
        };
      }
      default:
        return state;
    }
  }

  const [state, dispatch] = useReducer(reducer, initialState);

  function openModalReducer(state, action) {
    const { operation, target } = action.payload;
    // console.log("openModalState: ", state);
    // console.log("operation: ", operation, "target: ", target);
    switch (action.type) {
      case "set_modal_config": {
        return {
          ...state,
          operationType: operation,
          targetType: target,
        };
      }

      case "error": {
        console.log("error", action.payload);
        return {
          ...state,
          operationType: operation,
        };
      }
    }
  }

  const [openModalState, openModalDispatch] = useReducer(
    openModalReducer,
    initialOpenModalState,
  );

  useEffect(() => {
    const fetchData = async () => {
      // PreHomeへのリクエスト時
      if (location.pathname === locationPreHomePath) {
        // console.log("/prehome location: ", location);
        const response = await fetch(apiPreHomePath, {
          method: "GET",
          credentials: "same-origin",
        });
        setStatus(response.status);
        const json = await response.json();
        // もしもjson.Common.errorMsgが空の文字列じゃない場合、
        // 小さめのモーダルを表示させてjson.Common.errorMsgを表示し再ログインを促す
        if (json.Common.errorMsg !== "") {
          console.log("json.Common.errorMsg: ", json.Common.errorMsg);
          openModalDispatch({
            type: "error",
            payload: { operation: "error" },
          });
          setIsOpenModal(!isOpenModal);
          setErrors(json.Common.errorMsg);
          return;
        }

        // console.log("response json: ", json);
        setData({
          ...data,
          user_id: json.Common.user_id,
          username: json.Common.username,
          email: json.Common.email,
          phase: json.phase,
          projects: json.projects,
        });
      }

      if (location.pathname === locationHomePath) {
        const fetchDefaultData = async () => {
          const response = await fetch(apiHomePath, {
            method: "GET",
            credentials: "same-origin",
          });

          const json = await response.json();
          // もしもjson.Common.errorMsgが""じゃない場合、
          // 小さめのモーダルを表示させてjson.Common.errorMsgを表示し、最初の画面に戻るボタン以外の操作をさせないようにする。
          if (json.Common.errorMsg !== "") {
            console.log("json.Common.errorMsg: ", json.Common.errorMsg);
            openModalDispatch({
              type: "error",
              payload: { operation: "error" },
            });
            setIsOpenModal(!isOpenModal);
            setErrors(json.Common.errorMsg);
            return;
          }

          setData({
            ...data,
            user_id: json.Common.user_id,
            username: json.Common.username,
            email: json.Common.email,
            projects: json.projects,
            project: json.project,
            projectFolders: json.project_folders,
            filesPerFolder: json.files_per_folder,
          });
          setProjectPerFolder(
            fileAndFolderData(json.project_folders, json.files_per_folder),
          );
        };
        fetchDefaultData();
      }
    };

    fetchData();
  }, [location.pathname]);

  function switchToggleSidebar() {
    setToggleSidebar(!toggleSidebar);
  }

  const fileAndFolderData = (pf, fpf) => {
    // console.log("fileAndFolderData: ", pf, fpf);

    const folderMap = {};
    if (pf != null) {
      pf.forEach((folder) => {
        // console.log("folder確認: ", folder);
        folderMap[folder.folder_name] = {
          id: folder.folder_id,
          name: folder.folder_name,
          type: "folder",
          parent: folder.parent_id,
          children: [],
        };
      });

      // folder.childrenにparentがidと一致するものを入れる。
      Object.keys(folderMap).forEach((folderName) => {
        pf.forEach((item) => {
          if (folderMap[folderName].id === item.parent_id) {
            folderMap[folderName].children.unshift(folderMap[item.folder_name]);
          }
        });
      });
    }

    // 各folderに関連のあるfileを取り出してfolderMapのchildrenに入れる。
    if (fpf != null) {
      Object.keys(fpf).forEach((folderName) => {
        if (folderMap[folderName]) {
          // ↓folderが所持しているfileが入ったarray
          fpf[folderName].map((file) => {
            folderMap[folderName].children.push({
              id: file.file_id,
              name: file.filename,
              type: "file",
            });
          });
        }

        if (folderName === "null") {
          folderMap["null"] = {
            id: "",
            name: folderName,
            type: "null",
            parent: "",
            children: [],
          };
          fpf[folderName].map((file) => {
            folderMap["null"].children.push({
              id: file.file_id,
              name: file.filename,
              type: "file",
            });
          });
        }
      });
    }

    return Object.values(folderMap).filter((folder) => !folder.parent);
  };

  const NewFolderList = ({ foldersData, open }) => {
    // console.log("check open: ", open);
    // console.log("foldersData: ", foldersData);

    // folderをクリック時にselectFolderにidを保持させる。
    function holdClickedFolder(e) {
      const childrenElem = e.target.children;
      if (childrenElem[0] !== undefined) {
        setSelectFolder(childrenElem[0].getAttribute("data-folder-id"));
      } else {
        setSelectFolder(e.target.getAttribute("data-folder-id"));
      }
    }

    // fileクリック時にselectFileにidを保持させる。
    function holdClickedFile(e) {
      const childrenElem = e.target.children;
      if (childrenElem[0] !== undefined) {
        setSelectFile(childrenElem[0].getAttribute("data-file-id"));
      } else {
        setSelectFile(e.target.getAttribute("data-file-id"));
      }
    }

    function toggleFolder(name) {
      console.log("test toggleFolder: ", name);
      setOpenFolder((prev) => ({ ...prev, [name]: !prev[name] }));
    }

    function toggleFile(id, name) {
      // console.log("test toggleFile: ", name);

      const foldersFileList = Object.values(data.filesPerFolder);
      foldersFileList.forEach((files) => {
        files.forEach((file) => {
          if (file.file_id === id && file.filename === name) {
            dispatch({
              type: "open_file",
              payload: { id: id, filename: name },
            });
          }
        });
      });
    }

    return (
      <ul>
        {foldersData.map((item) => (
          <li
            key={item.name}
            className={"mx-5 mt-1 truncate text-[#9c3e78] "}
            // text-gray-600
            onMouseEnter={() => {}}
          >
            {item.type === "folder" && item.name !== "null" && (
              <div>
                <button
                  className={`flex justify-start w-full font-semibold hover:bg-[#e8e8e8] rounded-md ${item.id != null && selectFolder === item.id.toString() ? "bg-[#e8e8e8]" : "bg-[#fdfcff]"}`}
                  onClick={(event) => {
                    // console.log(
                    //   "確認: ",
                    //   "item.name: ",
                    //   item.name,
                    //   "item.type: ",
                    //   item.type,
                    //   "selectFolder: ",
                    //   item.id,
                    // );
                    toggleFolder(item.name);
                    holdClickedFolder(event);
                    openModalDispatch({
                      type: "set_modal_config",
                      payload: {
                        operation: openModalState.operationType,
                        target: "folder",
                      },
                    });
                  }}
                >
                  <span
                    className={"mr-auto pl-2"}
                    data-folder-id={item.id}
                    onClick={(event) => {
                      holdClickedFolder(event);
                    }}
                  >
                    {item.name}
                  </span>
                </button>
                {open[item.name] && (
                  <NewFolderList
                    foldersData={item.children}
                    open={openFolder}
                  />
                )}
              </div>
            )}
            {item.type === "file" && (
              <button
                className={`flex justify-start w-full hover:bg-[#6C235B] rounded-md ${item.id != null && selectFile === item.id.toString() ? "bg-[#6C235B] text-[#e8e8e8]" : "bg-[#fdfcff] text-gray-500"}`}
                onClick={(event) => {
                  toggleFile(item.id, item.name);
                  holdClickedFile(event);
                  openModalDispatch({
                    type: "set_modal_config",
                    payload: {
                      operation: openModalState.operationType,
                      target: "file",
                    },
                  });
                }}
              >
                <span
                  className={"pl-2 "}
                  data-file-id={item.id}
                  onClick={(event) => {
                    holdClickedFile(event);
                  }}
                >
                  {item.name}
                </span>
              </button>
            )}
            {item.type === "null" &&
              item.children.map((nullItem) => (
                <button
                  className={`flex justify-start w-full pl-2 hover:bg-[#6C235B] rounded-md ${nullItem.id != null && selectFile === nullItem.id.toString() ? "bg-[#6C235B] text-[#e8e8e8]" : "bg-[#fdfcff] text-gray-400"}`}
                  onClick={(event) => {
                    toggleFile(nullItem.id, nullItem.name);
                    holdClickedFile(event);
                    openModalDispatch({
                      type: "set_modal_config",
                      payload: {
                        operation: openModalState.operationType,
                        target: "file",
                      },
                    });
                  }}
                >
                  <span
                    className={"pl-2"}
                    data-file-id={nullItem.id}
                    onClick={(event) => {
                      holdClickedFile(event);
                    }}
                  >
                    {nullItem.name}
                  </span>
                </button>
              ))}
          </li>
        ))}
      </ul>
    );
  };

  NewFolderList.propTypes = {
    foldersData: PropTypes.array.isRequired,
    open: PropTypes.object.isRequired,
  };

  function handleMouseMove(e) {
    document.getSelection().removeAllRanges();

    console.log("e.clientX: ", e.clientX);
    if (e.clientX <= 125) {
      return setWidth(125);
    }

    if (e.clientX >= 600) {
      return setWidth(600);
    }

    setWidth(e.clientX);
  }

  function handleMouseUp() {
    const resizerElm = document.getElementsByClassName("resizer")[0];
    resizerElm.classList.remove("w-1");
    resizerElm.classList.add("w-[1px]");
    window.removeEventListener("mousemove", handleMouseMove);
  }

  function handleMouseDown() {
    const resizerElm = document.getElementsByClassName("resizer")[0];
    resizerElm.classList.remove("w-[1px]");
    resizerElm.classList.add("w-1");
    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);
  }

  function setDataFromJsonAfterCreateFolder(json) {
    if (json.Common.errorMsg !== "") {
      // console.log("json.Common.errorMsg: ", json.Common.errorMsg);
      openModalDispatch({
        type: "error",
        payload: { operation: "error" },
      });
      setIsOpenModal(!isOpenModal);
      setErrors(json.Common.errorMsg);
      return;
    }

    setData({
      ...data,
      projectFolders: json.project_folders,
    });
    setProjectPerFolder(
      fileAndFolderData(json.project_folders, json.files_per_folder),
    );
  }

  function setDataFromJsonAfterCreateFile(json) {
    if (json.Common.errorMsg !== "") {
      // console.log("json.Common.errorMsg: ", json.Common.errorMsg);
      openModalDispatch({
        type: "error",
        payload: { operation: "error" },
      });
      setIsOpenModal(!isOpenModal);
      setErrors(json.Common.errorMsg);
      return;
    }

    setData({
      ...data,
      filesPerFolder: json.files_per_folder,
    });
    setProjectPerFolder(
      fileAndFolderData(json.project_folders, json.files_per_folder),
    );
  }

  const NewFolderCreate = () => {
    return (
      <div className={"flex items-center justify-center  h-full"}>
        <div className={"px-5 w-full"}>
          <input
            id="new_folder_input"
            type="text"
            placeholder="新しく作成するフォルダ名を入力..."
            className={"w-full outline-none"}
            autoComplete={"off"}
            autoFocus={true}
            onBlur={() => {
              setIsOpenModal(!isOpenModal);
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                // console.log("new folder input Enter");
                const newFolderName =
                  document.getElementById("new_folder_input").value;
                const createNewFolder = async () => {
                  let parent_id = null;
                  if (selectFolder != null) {
                    parent_id = parseInt(selectFolder);
                  }

                  const response = await fetch(apiCreateNewFolderPath, {
                    method: "POST",
                    credentials: "same-origin",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                      user_id: data.user_id,
                      project_id: data.project.ID,
                      parent_id: parent_id,
                      folder_name: newFolderName,
                    }),
                  });

                  const json = await response.json();
                  setDataFromJsonAfterCreateFolder(json);
                };
                createNewFolder();
                setIsOpenModal(!isOpenModal);
              }
              if (e.key === "Escape") {
                setIsOpenModal(!isOpenModal);
              }
            }}
          />
        </div>
      </div>
    );
  };

  const NewFileCreate = () => {
    return (
      <div className={"flex items-center justify-center  h-full"}>
        <div className={"px-5 w-full"}>
          <input
            id="new_file_input"
            type="text"
            placeholder="新しく作成するファイル名を入力..."
            className={"w-full outline-none"}
            autoComplete={"off"}
            autoFocus={true}
            onBlur={() => {
              setIsOpenModal(!isOpenModal);
            }}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                // console.log("new file input Enter");
                const newFileName =
                  document.getElementById("new_file_input").value;

                const createNewFile = async () => {
                  let folder_id = null;
                  if (selectFolder != null) {
                    folder_id = parseInt(selectFolder);
                  }

                  const response = await fetch(apiCreateNewFilePath, {
                    method: "POST",
                    credentials: "same-origin",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                      user_id: data.user_id,
                      project_id: data.project.ID,
                      folder_id: folder_id,
                      filename: newFileName,
                    }),
                  });

                  const json = await response.json();
                  // console.log("api create new file response-json: ", json);
                  setDataFromJsonAfterCreateFile(json);
                  handleCreateFile(json.file.file_id, json.file.filename);
                };
                createNewFile();
                setIsOpenModal(!isOpenModal);
              }

              if (e.key === "Escape") {
                setIsOpenModal(!isOpenModal);
              }
            }}
          />
        </div>
      </div>
    );
  };

  function setDataFromJsonAfterDelete(json) {
    if (json.Common.errorMsg !== "") {
      console.log("json.Common.errorMsg: ", json.Common.errorMsg);
      openModalDispatch({
        type: "error",
        payload: { operation: "error" },
      });
      setIsOpenModal(!isOpenModal);
      setErrors(json.Common.errorMsg);
      return;
    }

    setData({
      ...data,
      projectFolders: json.project_folders,
      filesPerFolder: json.files_per_folder,
    });
    setProjectPerFolder(
      fileAndFolderData(json.project_folders, json.files_per_folder),
    );
  }

  function setDataFromJsonAfterDeleteFolder(json) {
    setDataFromJsonAfterDelete(json);

    // 削除された該当のfolderに紐づくfileをstate.openFile内から消去する
    Object.values(state.openFile).forEach((file, folder_id) => {
      if (folder_id === parseInt(file.folder_id)) {
        dispatch({
          type: "close_file",
          payload: { filename: file.filename },
        });
      }
    });
  }

  function setDataFromJsonAfterDeleteFile(json, filename) {
    setDataFromJsonAfterDelete(json);
    handleCloseFile(filename);
  }

  const DeleteFolder = () => {
    function getFolder() {
      let targetFolder;
      data.projectFolders.forEach((folder) => {
        if (folder.folder_id === parseInt(selectFolder)) {
          targetFolder = folder;
        }
      });
      return targetFolder;
    }
    const targetFolder = getFolder();

    const deleteFolder = async () => {
      let folder_id = null;
      if (selectFolder != null) {
        folder_id = parseInt(selectFolder);
      }

      const response = await fetch(apiDeleteFolderPath, {
        method: "DELETE",
        credentials: "same-origin",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          folder_id: folder_id,
          folder_name: targetFolder.folder_name,
          user_id: data.user_id,
          project_id: data.project.ID,
        }),
      });

      const json = await response.json();
      // console.log("api delete folder response-json: ", json);
      setDataFromJsonAfterDeleteFolder(json, folder_id);
    };

    return (
      <div className={"p-2"}>
        <div className={"flex justify-between"}>
          <h1 className={"text-xl font-bold"}>フォルダの削除</h1>
          <button
            className={"w-8 h-8 pr-1"}
            onClick={() => setIsOpenModal(!isOpenModal)}
          >
            <FontAwesomeIcon icon={faXmark} size={"lg"} color={"#6C235B"} />
          </button>
        </div>
        <br />

        <div className={"flex justify-center"}>
          <span>"{targetFolder.folder_name}" を削除しますか？</span>
        </div>
        <div className={"flex justify-center mt-8"}>
          <button
            className={
              "w-16 h-9 text-sm mr-3 bg-rose-200 text-rose-600 border-2 border-rose-600 rounded-md"
            }
            onClick={() => {
              deleteFolder();
              setIsOpenModal(!isOpenModal);
            }}
          >
            削除
          </button>
          <button
            className={
              "w-20 h-9 text-sm bg-gray-200 text-gray-700 border-2 border-gray-700 rounded-md"
            }
            onClick={() => setIsOpenModal(!isOpenModal)}
          >
            キャンセル
          </button>
        </div>
      </div>
    );
  };

  const DeleteFile = () => {
    // file_idは取得しているのでそれをもとにfileを取得する。
    function getFile() {
      let targetFile;
      Object.keys(data.filesPerFolder).forEach((folderName) => {
        data.filesPerFolder[folderName].forEach((file) => {
          // console.log("file: ", file);

          if (file.file_id === parseInt(selectFile)) {
            targetFile = file;
          }
        });
      });
      return targetFile;
    }

    const targetFile = getFile();

    const deleteFile = async () => {
      const response = await fetch(apiDeleteFilePath, {
        method: "DELETE",
        credentials: "same-origin",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          file_id: targetFile.file_id,
          filename: targetFile.filename,
          user_id: data.user_id,
          project_id: data.project.ID,
          folder_id: targetFile.parent_id,
        }),
      });

      const json = await response.json();
      setDataFromJsonAfterDeleteFile(json, targetFile.filename);
    };

    return (
      <div className={"p-2"}>
        <div className={"flex justify-between"}>
          <h1 className={"text-xl font-bold"}>ファイルの削除</h1>
          <button
            className={"w-8 h-8 pr-1"}
            onClick={() => setIsOpenModal(!isOpenModal)}
          >
            <FontAwesomeIcon icon={faXmark} size={"lg"} color={"#6C235B"} />
          </button>
        </div>
        <br />

        <div className={"flex justify-center"}>
          <span>"{targetFile.filename}" を削除しますか？</span>
        </div>
        <div className={"flex justify-center mt-8"}>
          <button
            className={
              "w-16 h-9 text-sm mr-3 bg-rose-200 text-rose-600 border-2 border-rose-600 rounded-md"
            }
            onClick={() => {
              deleteFile();
              setIsOpenModal(!isOpenModal);
            }}
          >
            削除
          </button>
          <button
            className={
              "w-20 h-9 text-sm bg-gray-200 text-gray-700 border-2 border-gray-700 rounded-md"
            }
            onClick={() => setIsOpenModal(!isOpenModal)}
          >
            キャンセル
          </button>
        </div>
      </div>
    );
  };

  function switchNoHighlightFolder(e) {
    const elem = e.target;
    const result = elem.classList.contains("side-bar-menu");
    if (selectFolder !== null && result) {
      setSelectFolder(null);
    }
  }

  const SideBarMenu = () => {
    // console.log("data.project: ", data.project);
    return (
      <>
        <div>
          <div className={"m-2"}>
            <div className={"relative top-0"}>
              <button
                className={"m-1"}
                onClick={() => {
                  openModalDispatch({
                    type: "set_modal_config",
                    payload: {
                      operation: "create",
                      target: "folder",
                    },
                  });
                  setIsOpenModal(!isOpenModal);
                }}
              >
                <FontAwesomeIcon
                  icon={faFolderPlus}
                  size={"xl"}
                  color={"#6C235B"}
                />
              </button>
              <button
                className={"m-1"}
                onClick={() => {
                  openModalDispatch({
                    type: "set_modal_config",
                    payload: {
                      operation: "create",
                      target: "file",
                    },
                  });
                  setIsOpenModal(!isOpenModal);
                }}
              >
                <FontAwesomeIcon
                  icon={faPenToSquare}
                  size={"xl"}
                  color={"#6C235B"}
                />
              </button>
              <button
                className={"m-1"}
                onClick={() => {
                  openModalDispatch({
                    type: "set_modal_config",
                    payload: {
                      operation: "delete",
                      target: openModalState.targetType,
                    },
                  });
                  setIsOpenModal(!isOpenModal);
                }}
              >
                <FontAwesomeIcon
                  icon={faTrashCan}
                  size={"xl"}
                  color={"#6C235B"}
                />
              </button>
            </div>
          </div>
          <NewFolderList foldersData={projectPerFolder} open={openFolder} />
          {/*<div className={"fixed bottom-1 mx-2"}>*/}
          {/*  <button className={"m-1"}>*/}
          {/*    <FontAwesomeIcon*/}
          {/*      icon={faCircleUser}*/}
          {/*      size={"2xl"}*/}
          {/*      color={"#6C235B"}*/}
          {/*    />*/}
          {/*  </button>*/}
          {/*</div>*/}
        </div>
      </>
    );
  };

  function handleCreateFile(id, filename) {
    dispatch({ type: "open_file", payload: { id: id, filename: filename } });
  }

  function handleCloseFile(filename) {
    dispatch({ type: "close_file", payload: { filename: filename } });
  }

  const OpenFileList = () => {
    // console.log("OpenFileList state.openFile : ", state.openFile);
    const list = Object.values(state.openFile);

    function switchFile(id, name) {
      // console.log("test: ", id, name, list);
      dispatch({ type: "open_file", payload: { id: id, filename: name } });
    }

    return (
      <ul className={"flex h-full list-none"}>
        {list.map((item) => {
          if (item) {
            return (
              <li
                key={item.file_id}
                itemID={item.file_id}
                className={
                  "flex-1 h-full min-w-0 text-center content-center border-1 border-gray-200 px-3 text-gray-400 hover:bg-[#6C235B] hover:text-[#fdfcff]"
                }
              >
                <div className={"flex items-center"}>
                  <div
                    className={"grow-1 text-left truncate"}
                    onClick={() => switchFile(item.file_id, item.filename)}
                  >
                    <span className={"block w-2/5 truncate"}>
                      {item.filename}
                    </span>
                  </div>
                  <div>
                    <button onClick={() => handleCloseFile(item.filename)}>
                      <FontAwesomeIcon
                        icon={faXmark}
                        size={"sm"}
                        color={"text-[#6C235B] hover:text-[#fdfcff]"}
                      />
                    </button>
                  </div>
                </div>
              </li>
            );
          }
        })}
      </ul>
    );
  };

  const ErrorDisplay = () => {
    console.log("errors: ", errors);
    return (
      <div className={"flex flex-col justify-center h-full"}>
        <div className={"flex flex-col justify-center items-center h-1/2"}>
          <span className={"font-bold text-rose-600"}>※{errors}</span>
          <span className={"text-rose-600"}>再度ログインしてください。</span>
        </div>
        <div className={"flex justify-center h-1/2"}>
          <button
            className="w-50 h-15 bg-[#6C235B] hover:bg-[#872C76] rounded text-white font-lg font-bold px-6 "
            type="button"
            onClick={() => {
              // console.log("再ログイン");
              navigate("/login");
            }}
          >
            Login
          </button>
        </div>
      </div>
    );
  };

  return (
    <>
      {location.pathname === locationPreHomePath ? (
        <Base>
          <PreHome projects={data.projects} />
        </Base>
      ) : (
        <>
          <div className={`wrapper h-screen bg-[#fdfcff]`}>
            {isOpenModal ? (
              <div
                className={
                  "overlay fixed inset-0 w-full h-full z-10 bg-gray-800 opacity-60"
                }
              ></div>
            ) : null}
            <div className={"header-wrapper"}>
              <header
                className="sidebar-header fixed top-0 h-10 z-1 bg-[#fdfcff]"
                style={{ width: width }}
              >
                <div className={"m-2"}>
                  <FontAwesomeIcon
                    icon={faArrowsLeftRight}
                    size={"xl"}
                    color={"#6C235B"}
                    className={"m-1"}
                    onClick={switchToggleSidebar}
                  />
                </div>
              </header>

              <header
                className={
                  "main-header fixed border-gray-200 border-1 border-b-0 top-0 h-10 "
                }
                style={{
                  left: width,
                  width: window.innerWidth - width,
                }}
              >
                <OpenFileList />
              </header>
            </div>

            {/*openModalState.operationTypeに"error"が入っている場合の処理*/}
            {openModalState.operationType === "error" && (
              <ErrorModal
                isOpenModal={isOpenModal}
                elements={<ErrorDisplay />}
              />
            )}

            {/*openModalState.targetTypeに値が入っていない場合の処理*/}
            {openModalState.operationType === "delete" &&
              openModalState.targetType === null &&
              null}

            {/*Modalをそれぞれの状況に合わせて表示する処理*/}
            {openModalState.operationType === "create" ? (
              openModalState.targetType === "file" ? (
                <FileFolderCreateModal
                  isOpenModal={isOpenModal}
                  elements={<NewFileCreate />}
                />
              ) : (
                <FileFolderCreateModal
                  isOpenModal={isOpenModal}
                  elements={<NewFolderCreate />}
                />
              )
            ) : openModalState.targetType ===
              null ? null : openModalState.targetType === "file" ? (
              <FileFolderDeleteModal
                isOpenModal={isOpenModal}
                elements={<DeleteFile />}
              />
            ) : (
              <FileFolderDeleteModal
                isOpenModal={isOpenModal}
                elements={<DeleteFolder />}
              />
            )}
            <div className="home flex h-full ">
              <div
                className={`${toggleSidebar ? "" : "-translate-x-full"} side-bar-menu mt-10 mb-5 border-y-1 border-gray-200  transform transition-transform duration-300 ease-in-out`}
                style={{ width: width }}
                onClick={switchNoHighlightFolder}
              >
                <SideBarMenu />
              </div>

              <div
                className={
                  "resizer relative w-[1px] h-full bg-gray-200 hover:bg-[#6C235B] cursor-col-resize"
                }
                onMouseDown={handleMouseDown}
                onMouseEnter={() => {
                  const resizerElm =
                    document.getElementsByClassName("resizer")[0];
                  resizerElm.classList.add("w-1");
                  resizerElm.classList.remove("w-[1px]");
                }}
                onMouseLeave={() => {
                  const resizerElm =
                    document.getElementsByClassName("resizer")[0];
                  resizerElm.classList.remove("w-1");
                  resizerElm.classList.add("w-[1px]");
                }}
              ></div>

              <div
                className={
                  "flex-grow mt-10 mb-5 border-1 border-l-0 border-gray-200 p-3"
                }
              >
                {Object.keys(state.openFile).length && state.currentFile ? (
                  <MilkdownProvider>
                    <HomeEditor
                      file={state.currentFile}
                      dispatch={dispatch}
                      data={data}
                      setData={setData}
                    />
                  </MilkdownProvider>
                ) : (
                  <div className={"flex justify-center text-gray-400"}>
                    <p>ファイルを選択してください</p>
                  </div>
                )}
              </div>
            </div>
          </div>
          <footer
            className={
              "home-footer fixed bottom-0 right-0 left-0 h-5 text-center bg-[#fdfcff]"
            }
          >
            <div className={"flex justify-center items-center"}>
              {state.currentFile ? (
                <span className={"font-extralight"}>
                  {state.currentFile.filename}
                </span>
              ) : (
                <span className={"font-extralight"}>dev-nikki</span>
              )}
            </div>
          </footer>
        </>
      )}
    </>
  );
};

export default Home;
