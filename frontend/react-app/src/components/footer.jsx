const Footer = () => {
  return (
    <>
      <footer className={"fixed bottom-0 right-0 left-0 text-center mb-7"}>
        {/*フォント変えるのありかもここだけ*/}
        <div className="bottom-proverb w-104 sm:w-112 md:w-120 lg:w-128 text-[24px] mx-auto mt-60">
          <h5 className="text-[#6C235B] font-semibold">To do is to be.</h5>
        </div>
        <p className="font-extralight">Jean-Paul Sartre</p>
      </footer>
    </>
  );
};

export default Footer;
