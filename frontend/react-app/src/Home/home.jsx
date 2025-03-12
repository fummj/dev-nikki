import { useEffect, useState } from "react";
import ErrorPage from "../components/error";

const Home = () => {
  const [data, setData] = useState({
    user_id: "",
    username: "",
    email: "",
    projects: [],
  });
  const [status, setStatus] = useState(0);

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("/api/home", {
        method: "GET",
        credentials: "same-origin",
      });
      setStatus(response.status);
      const json = await response.json();
      console.log("response json: ", json);
      setData({
        ...data,
        user_id: json.user_id,
        username: json.username,
        email: json.email,
      });
    };
    fetchData();
  }, []);

  if (400 <= status || 500 <= status) {
    return <ErrorPage code={status} />;
  }

  return (
    <>
      <p>status-code: {status}</p>
      <h1>Welcome Home</h1>
      <p>{data.user_id}</p>
      <p>{data.username}</p>
      <p>{data.email}</p>
    </>
  );
};

export default Home;
