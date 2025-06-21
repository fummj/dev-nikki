import { useMemo, useState } from "react";
import { useNavigate } from "react-router";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircleXmark } from "@fortawesome/free-solid-svg-icons";

const prehomePath = "/api/home/prehome";

const PreHome = ({ projects }) => {
  const [suggestState, setSuggestState] = useState([]);
  const [inputValue, setInputValue] = useState("");
  const [userKeyDownCountState, setUserKeyDownCountState] = useState(-1);
  const suggestCount = useMemo(() => suggestState.length, [suggestState]);

  const navigate = useNavigate();
  function handleHome() {
    navigate("/home");
  }

  const suggestList = suggestState.map((suggestion, index) => {
    return (
      <div key={index}>
        <ul
          itemID={index}
          className={
            "w-11/12 ml-5 pl-1 h-7 hover:bg-[#9D4889] hover:text-white text-left rounded-md outline-none"
          }
          onClick={handleClickSuggestion}
        >
          {suggestion}
        </ul>
      </div>
    );
  });

  function handleChangeSuggest(e, state, setState, setInputValue) {
    // 入力値をもとにprojects(props)の中身をsuggestState(リスト)に入れる。
    setState([]);
    setInputValue(e.target.value);

    let pattern = e.target.value.toLowerCase();
    for (let i = 0; i < state.length; i++) {
      let base = state[i].Name;
      if (base.startsWith(pattern) && pattern !== "") {
        setState((prev) => [...prev, state[i].Name]);
      }
    }

    if (e.target.value === "") {
      for (let i = 0; i < state.length; i++) {
        setState((prev) => [...prev, state[i].Name]);
      }
    }
  }

  function handleClickCloseButton(e, state, setState, setInputValue) {
    setState([]);
    setInputValue("");

    // state(projects)の中身をsuggestStateに入れる。
    for (let i = 0; i < state.length; i++) {
      setState((prev) => [...prev, state[i].Name]);
    }
  }

  function handleClickSuggestion(e) {
    let inputElm = document.getElementById("project_name_input");
    inputElm.value = e.target.innerHTML;
    inputElm.focus();
  }

  function handleKeyDown(
    e,
    countMemo,
    userKeyDownCount,
    setUserKeyDownCount,
    setInputValue,
  ) {
    let suggest_wrap = document.querySelector("#suggest-wrap");

    if (e.key === "Enter") {
      e.preventDefault();
      // console.log("project_name: ", inputValue);

      if (inputValue === "") {
        return;
      }

      const fetchData = async () => {
        const response = await fetch(prehomePath, {
          method: "POST",
          credentials: "same-origin",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ project_name: inputValue }),
        });
        const data = await response.json();
        if (data.Common.errorMsg === "") {
          // console.log("prehome response: ", data);
          // console.log("prehomeからhomeに遷移。");
          handleHome();
        }
      };
      fetchData();
    }

    if (e.key === "ArrowDown" && userKeyDownCount < countMemo - 1) {
      e.preventDefault();
      // console.log("userKeyDownCount: ", userKeyDownCount);
      handleArrowDown(
        suggest_wrap,
        userKeyDownCount,
        setUserKeyDownCount,
        setInputValue,
        countMemo,
      );
    }

    if (e.key === "ArrowUp" && userKeyDownCount > 0) {
      e.preventDefault();
      handleArrowUp(
        suggest_wrap,
        userKeyDownCount,
        setUserKeyDownCount,
        setInputValue,
      );
    }

    if (e.ctrlKey) {
      if (e.key === "n" && userKeyDownCount < countMemo - 1) {
        e.preventDefault();
        // console.log("same ↓");
        handleArrowDown(
          suggest_wrap,
          userKeyDownCount,
          setUserKeyDownCount,
          setInputValue,
          countMemo,
        );
      }

      if (e.key === "p" && userKeyDownCount > 0) {
        e.preventDefault();
        // console.log("same ↑");
        handleArrowUp(
          suggest_wrap,
          userKeyDownCount,
          setUserKeyDownCount,
          setInputValue,
        );
      }
    }
  }

  function handleArrowDown(
    elm,
    keyDownCount,
    setKeyDownCount,
    setInputValue,
    memo,
  ) {
    // 前回の要素の色変更を解除
    if (keyDownCount >= 0 && keyDownCount < memo - 1) {
      let prevSuggestElm = elm.querySelector(`[itemid="${keyDownCount}"]`);
      prevSuggestElm.classList.remove("bg-[#9D4889]", "text-white");
    }

    let count = keyDownCount + 1;
    let suggestElm = elm.querySelector(`[itemid="${count}"]`);

    // console.log(suggestElm);
    suggestElm.classList.add("bg-[#9D4889]", "text-white");

    setInputValue(suggestElm.innerHTML);

    setKeyDownCount(keyDownCount + 1);
  }

  function handleArrowUp(elm, keyDownCount, setKeyDownCount, setInputValue) {
    // 前回の要素の色変更を解除
    let prevSuggestElm = elm.querySelector(`[itemid="${keyDownCount}"]`);
    prevSuggestElm.classList.remove("bg-[#9D4889]", "text-white");

    let count = keyDownCount - 1;
    let suggestElm = elm.querySelector(`[itemid="${count}"]`);

    suggestElm.classList.add("bg-[#9D4889]", "text-white");

    setInputValue(suggestElm.innerHTML);

    setKeyDownCount(keyDownCount - 1);
  }

  return (
    <>
      <div className={"w-112 m-auto"}>
        <form
          className={
            "sm:w-96 md:w-104 lg:w-112 h-92 overflow-y-auto no_scrollbar border-7 border-[#6C235B] outline-none rounded-3xl font-extralight"
          }
          id="form"
          action={prehomePath}
          method="POST"
        >
          <div className={"mx-auto w-11/12 flex h-10 m-1"}>
            <input
              className={"w-full outline-none rounded-sm m-1 px-1"}
              id="project_name_input"
              type="text"
              name="project_name"
              value={inputValue}
              placeholder="Project Name"
              required={true}
              autoComplete={"off"}
              autoFocus={true}
              onChange={(e) =>
                handleChangeSuggest(e, projects, setSuggestState, setInputValue)
              }
              onKeyDown={(e) =>
                handleKeyDown(
                  e,
                  suggestCount,
                  userKeyDownCountState,
                  setUserKeyDownCountState,
                  setInputValue,
                )
              }
            />
            <button
              className={" mr-1"}
              type="reset"
              onClick={(e) =>
                handleClickCloseButton(
                  e,
                  projects,
                  setSuggestState,
                  setInputValue,
                )
              }
            >
              <FontAwesomeIcon
                icon={faCircleXmark}
                size={"lg"}
                color={"#6C235B"}
              />
            </button>
          </div>
          <hr className={"mx-auto w-11/12 text-gray-200 border-b-1"} />
          <div
            className={"flex flex-col mt-3 align-center gap-1 "}
            id={"suggest-wrap"}
          >
            {suggestList}
          </div>
        </form>
      </div>
    </>
  );
};

export default PreHome;
