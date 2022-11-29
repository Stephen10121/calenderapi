import './App.css';
import { getCookie } from "../functions/cookie";
import { useEffect, useState } from 'react';
import NotLogged from './components/NotLogged';
import { validate } from "../functions/backendFetch";
import LoggedIn from './components/LoggedIn';

export default function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [token, setToken] = useState("");

  useEffect(() => {
    const cookie = getCookie("G_CAL");
    if (cookie) {
      validate(cookie).then((data) => {
        if (data.error || !data.data) {
          return;
        }
        setName(data.data.name);
        setEmail(data.data.email);
        setToken(cookie);
        setLoggedIn(true);
      });
    }
  }, []);

  function logout() {
    document.cookie = "G_CAL=";
    setLoggedIn(false);
    setName("");
    setEmail("");
    setToken("");
  }

  function isLoggedIn(email: string, name: string, token: string) {
    setName(name);
    setEmail(email);
    setToken(token);
    setLoggedIn(true);
  }

  if (loggedIn) {
    return <LoggedIn name={name} email={email} token={token} logout={logout}/>
  }
  return <NotLogged loggedIn={isLoggedIn}/>;
}