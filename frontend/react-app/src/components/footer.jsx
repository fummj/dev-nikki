const Footer = () => {
  return (
    <>
      <footer className={"fixed bottom-0 right-0 left-0 text-center mb-7"}>
        {/*フォント変えるのありかもここだけ*/}
        <div className={"flex flex-col justify-center"}>
          <div className="bottom-proverb text-[24px] ">
            <h5 className="text-[#6C235B] font-semibold">To do is to be.</h5>
          </div>
          <p className="font-extralight">Jean-Paul Sartre</p>
        </div>
      </footer>
    </>
  );
};

export default Footer;
