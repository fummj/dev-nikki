import { useNavigate } from "react-router";

const HeaderNav = () => {
  function switchShowButton() {
    let path = window.location.pathname;
    if (path !== "/") {
      console.log("not /");
      return "hidden";
    } else {
      console.log("yes /");
      return "inline-block";
    }
  }

  const navigate = useNavigate();

  function handleRoot() {
    console.log("index page");
    navigate("/");
  }

  function handleLogin() {
    console.log("login page");
    navigate("/login");
  }

  function handleSignUp() {
    console.log("sign-up page");
    navigate("/signup");
  }

  return (
    <header>
      <div className="header-inner mt-5 mr-8">
        <nav className="flex items-center justify-between">
          <div className="header-logo" onClick={handleRoot}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="160"
              height="80"
              viewBox="0 0 100 100"
            >
              <path d="M78.001 44.011c-.399-.23-.729-.798-.729-1.261V26.522c0-.461-.326-1.028-.726-1.258l-16.965-9.795c-.399-.231-1.053-.231-1.454 0l-14.054 8.114c-.398.23-1.053.23-1.453 0l-14.056-8.114c-.397-.23-1.053-.23-1.453 0l-16.965 9.795c-.398.23-.727.797-.727 1.258v19.59c0 .462.328 1.03.727 1.26l14.056 8.114c.4.23.728.798.728 1.259v16.227c0 .461.326 1.029.727 1.259l16.964 9.796c.4.23 1.055.23 1.453 0l14.056-8.115c.399-.229 1.054-.229 1.454 0l14.054 8.115c.399.23 1.054.23 1.454 0l16.963-9.796c.4-.23.727-.798.727-1.259V53.384c0-.462-.326-1.029-.727-1.26l-14.054-8.113zm-18.417 6.433c-.4.231-1.055.231-1.454 0l-8.232-4.753c-.398-.23-.398-.609 0-.839l8.236-4.757c.398-.23 1.053-.23 1.454 0l8.228 4.755c.399.231.399.609 0 .84l-8.232 4.754zm-16.963-1.393c.4-.23 1.055-.23 1.453.002l8.233 4.751c.399.232.399.61 0 .84L44.07 59.4c-.398.23-1.054.23-1.455 0l-8.227-4.754c-.4-.231-.4-.609 0-.841l8.233-4.754zm-15.509-26.86c.4-.23 1.056-.23 1.453 0l8.234 4.752c.398.231.398.61 0 .841l-8.236 4.754c-.399.231-1.054.231-1.453 0l-8.23-4.753c-.399-.231-.399-.609 0-.84l8.232-4.754zM15.241 33.248c0-.462.328-.651.729-.42l11.138 6.434c.4.231 1.055.231 1.454.001l29.566-17.072c.401-.23 1.055-.23 1.454 0l11.145 6.434c.397.229.726.796.726 1.259v9.503c0 .461-.327.651-.726.419l-11.139-6.434c-.401-.231-1.055-.231-1.454-.001L28.565 50.444c-.397.231-1.053.231-1.453 0L15.97 44.011c-.4-.23-.729-.797-.729-1.26v-9.503zm16.238 37.624c-.4-.23-.727-.797-.727-1.261v-9.503c0-.462.326-.65.727-.418l11.137 6.433c.401.232 1.057.232 1.455 0l29.567-17.072c.399-.23 1.054-.23 1.454 0l11.143 6.433c.4.23.727.798.727 1.26v9.505c0 .461-.326.65-.727.419l-11.139-6.434c-.4-.231-1.054-.231-1.453-.001L44.074 77.305c-.398.231-1.053.231-1.453 0l-11.142-6.433zm33.924 1.68c-.401-.23-.401-.608 0-.838l8.237-4.758c.401-.23 1.055-.23 1.455.001l8.229 4.754c.401.23.401.609 0 .84l-8.233 4.754c-.4.231-1.055.231-1.454 0l-8.234-4.753z" />
            </svg>
          </div>
          <div className={`switchShowButton() flex`}>
            <button
              className="bg-white hover:bg-gray-100 border border-gray-500 rounded text-gray-900  py-3 px-5 m-2"
              type="button"
              onClick={handleLogin}
            >
              Login
            </button>
            <button
              className="bg-[#6C235B] hover:bg-[#872C76] rounded text-white  py-3 px-3 m-2"
              type="button"
              onClick={handleSignUp}
            >
              Sign Up
            </button>
          </div>
        </nav>
      </div>
    </header>
  );
};

export default HeaderNav;
