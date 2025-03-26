import { useEffect, useState } from "react";
import { useLocation } from "react-router";
import ErrorPage from "../components/error";
import PreHome from "./prehome.jsx";
import Base from "../Index/base.jsx";

const locationPreHomePath = "/prehome";
const locationHomePath = "/home";
const apiPreHomePath = "/api/home/prehome";

const Home = () => {
  const [data, setData] = useState({
    user_id: "",
    username: "",
    email: "",
    errMsg: "",
    phase: "",
    projects: [],
    project: {},
    projectFolders: [],
    filesPerFolder: {},
  });
  const [status, setStatus] = useState(0);
  const location = useLocation();

  useEffect(() => {
    const fetchData = async () => {
      if (location.pathname === locationPreHomePath) {
        console.log("/prehome location: ", location);
        const response = await fetch(apiPreHomePath, {
          method: "GET",
          credentials: "same-origin",
        });
        setStatus(response.status);
        console.log("response: ", response);
        const json = await response.json();
        console.log("response json: ", json);
        setData({
          ...data,
          user_id: json.Common.user_id,
          username: json.Common.username,
          email: json.Common.email,
          phase: json.phase,
          projects: json.projects,
        });
      }

      if (location.pathname === locationHomePath) {
        console.log("/home location.state: ", location.state);
        setData({
          user_id: location.state.userData.user_id,
          username: location.state.userData.username,
          email: location.state.userData.email,
          errMsg: location.state.userData.errMsg,
          phase: location.state.userData.phase,
          projects: location.state.userData.projects,
          project: location.state.userData.project,
          projectFolders: location.state.userData.projectFolders,
          filesPerFolder: location.state.userData.filesPerFolder,
        });
      }
    };

    fetchData();
  }, [location]);

  if (400 <= status || 500 <= status) {
    return <ErrorPage code={status} />;
  }

  return (
    <>
      {location.pathname === locationPreHomePath ? (
        <Base>
          <PreHome projects={data.projects} />
        </Base>
      ) : (
        <>
          <p>status: {status}</p>
          <h1>Welcome Home</h1>
          <p>{data.user_id}</p>
          <p>{data.username}</p>
          <p>{data.email}</p>
        </>
      )}
    </>
  );
};

export default Home;
