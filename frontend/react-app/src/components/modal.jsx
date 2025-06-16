const Modal = ({ isOpenModal, elements }) => {
  if (!isOpenModal) {
    return null;
  }

  return (
    <div
      className={
        "fixed z-100 inset-0 m-auto mt-20 sm:w-112 md:w-128 lg:w-156 h-20 bg-[#fdfcff] border-5 border-[#6C235B] rounded-xl"
      }
    >
      {elements}
    </div>
  );
};

const FileFolderCreateModal = ({ isOpenModal, elements }) => {
  if (!isOpenModal) {
    return null;
  }

  return (
    <div
      className={
        "fixed z-100 inset-0 m-auto mt-20 sm:w-112 md:w-128 lg:w-156 h-20 bg-[#fdfcff] border-5 border-[#6C235B] rounded-xl"
      }
    >
      {elements}
    </div>
  );
};

const FileFolderDeleteModal = ({ isOpenModal, elements }) => {
  if (!isOpenModal) {
    return null;
  }

  return (
    <div
      className={
        "fixed z-100 inset-0 m-auto w-90 h-45 bg-[#fdfcff] border-5 border-[#6C235B] rounded-xl"
      }
    >
      {elements}
    </div>
  );
};

const ErrorModal = ({ isOpenModal, elements }) => {
  if (!isOpenModal) {
    return null;
  }

  return (
    <div
      className={
        "fixed z-100 inset-0 m-auto w-90 h-45 bg-[#fdfcff] border-5 border-[#6C235B] rounded-xl"
      }
    >
      {elements}
    </div>
  );
};

export default Modal;
export { ErrorModal, FileFolderCreateModal, FileFolderDeleteModal };
