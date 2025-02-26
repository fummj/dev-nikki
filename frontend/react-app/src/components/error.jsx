const statusCodes = {
  401: "Unauthorized",
  403: "Forbidden",
  404: "Not Found",
  500: "Internal Server Error",
};

const ErrorPage = ({ code }) => {
  return (
    <>
      <h1>
        {code}
        {statusCodes[code]}
      </h1>
    </>
  );
};

export default ErrorPage;
