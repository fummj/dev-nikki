import { SignUpForm } from "./../components/form.jsx";
import Base from "../Index/base.jsx";

const SignUpElements = () => {
  return (
    <>
      <div className="flex flex-col items-center gap-5">
        <h2 className={"text-[32px] font-bold text-[#6C235B]"}>Sign Up</h2>
        <SignUpForm />
      </div>
    </>
  );
};

const SignUp = () => {
  return <Base>{<SignUpElements />}</Base>;
};

export default SignUp;
