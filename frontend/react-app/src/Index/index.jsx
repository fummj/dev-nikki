import HeaderNav from "./../components/header_nav.jsx";
import "./index.css";
import { useNavigate } from "react-router";

const Index = () => {
  const navigate = useNavigate();
  function handleSignUp() {
    console.log("sign-up page");
    navigate("/signup");
  }
  return (
    <>
      <HeaderNav />
      <main className="text-center mx-10">
        <div className="index-container">
          <div className="w-104 sm:w-112 md:w-120 lg:w-128 mx-auto">
            <img className="base-logo" src="/dev_nikki.png" alt="dev-nikki" />
          </div>
          <div className="mx-auto">
            <h2 className="index-heading text-[64px] font-bold mb-5">
              開発の記録を共有する
            </h2>
            <span className="w-52 sm:w-56 md:w-60 lg:w-64 text-base  text-gray-500 ">
              dev-nikkiは何かを個人で開発されている方に向けたサービスです。
              <br />
              その日に行った開発内容を日記として記録し、後から自身で振り返ることができます。
              <br />
              また、開発記録を他のユーザーと共有することもできます。
            </span>
          </div>
          <div className="w-104 sm:w-112 md:w-120 lg:w-128 mt-15 mb-20 mx-auto">
            <button
              className="bg-[#6C235B] hover:bg-[#872C76] rounded text-white font-bold py-3 px-6 "
              type="button"
              onClick={handleSignUp}
            >
              始める
            </button>
          </div>
        </div>
        {/*フォント変えるのありかもここだけ*/}
        <div className="bottom-proverb w-104 sm:w-112 md:w-120 lg:w-128 text-[24px] font-bold mx-auto mt-60">
          <h5 className="text-[#6C235B]">To do is to be.</h5>
        </div>
      </main>

      <div id="react-form"></div>
    </>
  );
};

export default Index;
