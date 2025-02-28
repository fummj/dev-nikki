import HeaderNav from "../components/header_nav.jsx";

const Base = (props) => {
  return (
    <>
      <HeaderNav />
      <main className="text-center mx-10">
        <div className="base-container">
          <div className="w-104 sm:w-112 md:w-120 lg:w-128 mx-auto">
            <img className="base-logo" src="/dev_nikki.png" alt="dev-nikki" />
          </div>
          {props.children}
        </div>
      </main>
      <footer className="text-center mx-10">
        {/*フォント変えるのありかもここだけ*/}
        <div className="bottom-proverb w-104 sm:w-112 md:w-120 lg:w-128 text-[24px] font-bold mx-auto mt-60">
          <h5 className="text-[#6C235B]">To do is to be.</h5>
        </div>
      </footer>
    </>
  );
};

export default Base;
