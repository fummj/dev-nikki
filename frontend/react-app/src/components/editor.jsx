import { useEffect, useRef } from "react";
import { Crepe } from "@milkdown/crepe";

import "@milkdown/crepe/theme/common/style.css";
import "@milkdown/crepe/theme/nord.css";
import "./editor.css";

const updateMarkdownPath = "/api/home/updateMarkdown";

const DEBOUNCE_DELAY = 5000;
let debounceTimer;

const HomeEditor = ({ file, dispatch, data, setData }) => {
  const editorRef = useRef(null);
  const editorInstanceRef = useRef(null);

  // console.log("rendering file: ", file);

  const sendMarkdownContent = (markdown) => {
    // console.log("5 seconds later sendMarkdownContent", markdown);

    const updateMarkdown = async () => {
      // console.log("file: ", file);

      const response = await fetch(updateMarkdownPath, {
        method: "PUT",
        credentials: "same-origin",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: file.user_id,
          project_id: file.project_id,
          folder_id: file.folder_id,
          file_id: file.file_id,
          filename: file.filename,
          content: markdown,
        }),
      });
      const json = await response.json();
      // console.log("json.file: ", json.file);
      dispatch({ type: "update_file", payload: { file: json.file } });
      setData({
        ...data,
        filesPerFolder: json.files_per_folder,
      });
    };
    updateMarkdown();
  };

  const debouncedSendMarkdownContent = (markdown) => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      sendMarkdownContent(markdown);
    }, DEBOUNCE_DELAY);
  };

  useEffect(() => {
    if (!editorRef.current) return;

    editorRef.current.innerHTML = "";

    if (editorInstanceRef.current) {
      editorInstanceRef.current.destroy();
      editorInstanceRef.current = null;
    }

    const Editor = new Crepe({
      root: editorRef.current,
      defaultValue: file.content,
    });

    // console.log("Editor: ", Editor);

    Editor.on((listener) => {
      listener.markdownUpdated((ctx, markdown, preMarkdown) => {
        debouncedSendMarkdownContent(markdown);
      });
    });

    Editor.create().then(() => {
      const markdownContent = Editor.getMarkdown();
      // console.log("markdownContent", markdownContent);
    });

    editorInstanceRef.current = Editor;
    // console.log("確認1 editorInstanceRef", editorInstanceRef.current);
  }, [file]);

  return <div ref={editorRef}></div>;
};

export default HomeEditor;
