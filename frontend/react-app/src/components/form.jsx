import { useState } from "react";
import { useNavigate } from "react-router";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEyeSlash } from "@fortawesome/free-solid-svg-icons";
import { faEye } from "@fortawesome/free-regular-svg-icons";

import "./form.css";

const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
const passwordPattern =
  /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
const loginPath = "/api/login";
const signupPath = "/api/signup";
const oauth2Path = "/api/auth/login";
const prehomePath = "/prehome";

const LoginForm = (isLogin) => {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [error, setError] = useState({
    email: "",
    password: "",
    responseMsg: "",
  });
  const [isRevealPassword, setIsRevealPassword] = useState(false);

  const navigate = useNavigate();
  function handlePreHome() {
    navigate(prehomePath);
  }

  function handleOAuth2() {
    console.log("OAuth2");
    window.location.href = oauth2Path;
  }

  function emailValidation(event) {
    setFormData({ ...formData, email: event.target.value });
    if (emailPattern.test(event.target.value)) {
      console.log("valid email");
      setError({ ...error, email: "" });
    } else {
      console.log("invalid email");
      setError({
        ...error,
        email: "※無効なメールアドレスです。",
      });
    }
  }

  function passwordValidation(event) {
    setFormData({ ...formData, password: event.target.value });
    if (passwordPattern.test(event.target.value)) {
      console.log("valid password");
      setError({ ...error, password: "" });
    } else {
      console.log("invalid password");
      setError({
        ...error,
        password:
          "※パスワードは8文字以上で、「大文字」「小文字」「数字」「記号」をそれぞれ1つ以上含めてください。",
      });
    }
  }

  const togglePassword = () => {
    setIsRevealPassword((isRevealPassword) => !isRevealPassword);
  };

  async function fetchResultAuth(isLogin, navFunc) {
    let formData = new FormData(document.getElementById("form"));
    console.log("formData: ", formData);
    try {
      const response = await fetch(isLogin ? loginPath : signupPath, {
        method: "POST",
        body: formData,
      });
      const data = await response.json();
      if (data.Common.errorMsg !== "") {
        setError({
          ...error,
          responseMsg: data.Common.errorMsg,
        });
        console.log(data);
        console.log("login failed: ", data.Common.errorMsg, error);
      } else {
        navFunc();
      }
      console.log(data);
    } catch (error) {
      console.log("Error fetch data: ", error);
    }
  }

  const NameInput = () => {
    return (
      <input
        className={
          "w-64 sm:w-72 md:w-80 lg:w-104 h-12 border-4 border-[#6C235B] outline-none rounded py-2 px-4"
        }
        type="text"
        name="name"
        id="name"
        placeholder="名前"
        required
      />
    );
  };

  const displayErrorMsg = () => {
    if (error.email !== "" && error.password !== "") {
      return (
        <>
          <span className={"flex flex-col"}>
            {error.email}
            <br />
            {error.password}
          </span>
        </>
      );
    } else if (error.responseMsg !== "") {
      return (
        <div>
          <span className={"block"}>{error.responseMsg}</span>
          <span>
            {error.email}
            {error.password}
          </span>
        </div>
      );
    } else {
      return (
        <span>
          {error.email}
          {error.password}
          {error.responseMsg}
        </span>
      );
    }
  };

  return (
    <>
      <div className={"w-104 m-auto"}>
        <form
          id="form"
          action=""
          method="POST"
          className={"flex flex-col justify-center items-center gap-3"}
        >
          <button
            type="button"
            className={
              "flex justify-start items-center w-64 sm:w-72 md:w-80 lg:w-104 h-12 border-2 border-gray-300 outline-none rounded py-2 px-4"
            }
            onClick={handleOAuth2}
          >
            <img
              src="/google_logo.png"
              alt="google-logo"
              width={"25px"}
              height={"10px"}
            />
            <span className={"w-full pr-10"}>Googleで続ける</span>
          </button>
          <div
            className={"flex justify-center w-64 sm:w-72 md:w-80 lg:w-104 py-4"}
          >
            <hr className={"w-full text-gray-300 border-b-0"} />
          </div>
          {isLogin ? console.log("LoginForm") : NameInput()}
          <input
            className={
              "w-64 sm:w-72 md:w-80 lg:w-104 h-12 border-4 border-[#6C235B] outline-none rounded py-2 px-4"
            }
            type="text"
            name="email"
            id="email"
            placeholder="メールアドレス"
            required
            onChange={(e) => {
              emailValidation(e);
            }}
          />
          <div
            className={
              "flex justify-between w-64 sm:w-72 md:w-80 lg:w-104 h-12  border-4 border-[#6C235B]  rounded "
            }
          >
            <input
              className={"w-full outline-none py-2 pl-4 pr-2 pt-4 pb-4"}
              type={isRevealPassword ? "text" : "password"}
              name="password"
              id="password"
              placeholder="パスワード"
              required
              onChange={(e) => {
                passwordValidation(e);
              }}
            />
            <span
              className={"py-2 px-2 pl-0"}
              onClick={() =>
                togglePassword(isRevealPassword, setIsRevealPassword)
              }
            >
              <FontAwesomeIcon icon={isRevealPassword ? faEye : faEyeSlash} />
            </span>
          </div>
          <input
            className={
              "w-64 sm:w-72 md:w-80 lg:w-104 h-12 bg-[#6C235B] hover:border[#994a7b] hover:bg-[#994a7b] disabled:opacity-70 disabled:cursor-not-allowed rounded text-white text-center"
            }
            type="button"
            value="送信"
            disabled={error.email !== "" || error.password !== ""}
            onClick={() => {
              if (error.email === "" && error.password === "") {
                fetchResultAuth(isLogin, handlePreHome);
              }
              console.log(error);
            }}
          />
          <div className={"flex flex-col text-sm text-center text-rose-600"}>
            {/*全体のエラーメッセージを表示*/}
            {displayErrorMsg(error)}
          </div>
        </form>
      </div>
    </>
  );
};

const SignUpForm = () => {
  let isLogin = false;
  return LoginForm(isLogin);
};

export { LoginForm, SignUpForm };
