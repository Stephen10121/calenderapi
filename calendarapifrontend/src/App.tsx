import './App.css';
import { getCookie } from "../functions/cookie";
import { POST_SERVER } from "../functions/variables";
import { useState } from 'react';

async function fetchGroups(token: string) {
  try {
    const groups = await fetch(`${POST_SERVER}/myGroups`, {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${token}`
      },
      credentials: "omit"
    })
    console.log(groups);
  } catch (err) {
    console.error(err);
    return false;
  }
}

function App() {
  const cookie = getCookie("G_CAL");
  const [loggedIn, setLoggedIn] = useState(cookie ? true : false);

  if (cookie) {
    fetchGroups(cookie).then((data) => {
      if (!data) {
        setLoggedIn(false);
      }
    });
  }
  return (
    <div className="App">
      {loggedIn ? "logged in" : "not logged in"}
    </div>
  );
}

export default App;
