import { LoginForm } from "./../components/form.jsx";
import Base from "../Index/base.jsx";

const LoginElements = () => {
  return (
    <>
      <div className="flex flex-col justify-center items-center gap-5">
        <h2 className={"text-[32px] font-bold text-[#6C235B]"}>Login</h2>
        <LoginForm isLogin={true} />
      </div>
    </>
  );
};

const Login = () => {
  return <Base>{<LoginElements />}</Base>;
};

export default Login;
