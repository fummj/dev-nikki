import { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEyeSlash } from "@fortawesome/free-solid-svg-icons";
import { faEye } from "@fortawesome/free-regular-svg-icons";

import HeaderNav from "./header_nav.jsx";
import "./form.css";

const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
const passwordPattern =
  /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
const loginPath = "/api/login";
const signupPath = "/api/signup";

function emailValidation(
  event,
  state,
  setEmailStateFunc,
  errorState,
  setErrorStateFunc,
) {
  setEmailStateFunc({ ...state, email: event.target.value });
  if (emailPattern.test(event.target.value)) {
    console.log("valid email");
    setErrorStateFunc({ ...errorState, email: "" });
  } else {
    console.log("invalid email");
    setErrorStateFunc({ ...errorState, email: "※無効なメールアドレスです。" });
  }
}

function passwordValidation(
  event,
  state,
  setPasswordStateFunc,
  errorState,
  setErrorStateFunc,
) {
  setPasswordStateFunc({ ...state, password: event.target.value });
  if (passwordPattern.test(event.target.value)) {
    console.log("valid password");
    setErrorStateFunc({ ...errorState, password: "" });
  } else {
    console.log("invalid password");
    setErrorStateFunc({
      ...errorState,
      password:
        "※パスワードは8文字以上で、「大文字」「小文字」「数字」「記号」をそれぞれ1つ以上含めてください。",
    });
  }
}

async function getData(isLogin) {
  let formData = new FormData(document.getElementById("form"));
  console.log("formData: ", formData);
  // const response = await fetch("http://localhost:8080/api/login", {
  const response = await fetch(isLogin ? loginPath : signupPath, {
    method: "POST",
    body: formData,
  });
  const data = await response.json();
  console.log(data);
}

const togglePassword = (isRevealState, setIsRevealStateFunc) => {
  setIsRevealStateFunc((isRevealState) => !isRevealState);
};

const displayError = (errorState) => {
  if (errorState.email !== "" && errorState.password !== "") {
    return (
      <>
        <span className={"flex flex-col"}>
          {errorState.email}
          <br />
          {errorState.password}
        </span>
      </>
    );
  } else {
    return (
      <span>
        {errorState.email}
        {errorState.password}
      </span>
    );
  }
};

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

const LoginForm = (isLogin) => {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [error, setError] = useState({
    email: "",
    password: "",
  });
  const [isRevealPassword, setIsRevealPassword] = useState(false);

  return (
    <>
      <HeaderNav />
      <div className={"w-104 m-auto"}>
        <form
          id="form"
          action=""
          method="POST"
          className={"flex flex-col justify-center items-center gap-2"}
        >
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
            onChange={(e) =>
              emailValidation(e, formData, setFormData, error, setError)
            }
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
              onChange={(e) =>
                passwordValidation(e, formData, setFormData, error, setError)
              }
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
              "w-64 sm:w-72 md:w-80 lg:w-104 h-12   bg-[#6C235B] hover:border[#994a7b] hover:bg-[#994a7b] rounded text-white text-center"
            }
            type="button"
            value="送信"
            onClick={() => getData(isLogin)}
          />
          <div className={"flex flex-col text-sm text-center text-rose-600"}>
            {/*email, passwordのエラーメッセージを表示*/}
            {displayError(error)}
          </div>
        </form>
        {/*確認用*/}
        {/*<p aria-disabled={true}>Form Data: {JSON.stringify(formData)}</p>*/}
      </div>
    </>
  );
};

const SignUpForm = () => {
  let isLogin = false;
  return LoginForm(isLogin);
};

export { LoginForm, SignUpForm };
