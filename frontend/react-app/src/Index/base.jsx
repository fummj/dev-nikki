import HeaderNav from "../components/header_nav.jsx";
import Footer from "../components/footer.jsx";

const Base = (props) => {
  return (
    <>
      <HeaderNav />
      <main className="text-center mx-10">
        <div className="base-container">
          <div className="w-80 sm:w-88 md:w-96 lg:w-104 mx-auto">
            <img className="base-logo" src="/dev_nikki.png" alt="dev-nikki" />
          </div>
          {props.children}
        </div>
      </main>
      <Footer />
    </>
  );
};

export default Base;
