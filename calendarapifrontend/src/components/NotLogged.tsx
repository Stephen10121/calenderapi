import { FormEvent, FormEventHandler, useState } from "react";
import { login } from "../../functions/backendFetch";
import classes from "./NotLogged.module.css";

interface ResponseData {
  email: string
  name: string
  token: string
}

export default function NotLogged({ loggedIn }: { loggedIn: (email: string, password: string, token: string) => void }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  function formSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    login(username, password).then((data) => {
      if (data.error) {
        setError(data.error);
        return
      }
      const response = data.data as ResponseData;
      document.cookie = `G_CAL=${response.token}`;
      loggedIn(response.email, response.name, response.token);
    });
  }

  return (
    <div className={classes.notLogged}>
      <h1 className={classes.welcome}>Welcome</h1>
      {error.length !== 0 ? error : null}
      <form onSubmit={formSubmit}>
        <input className={classes.input} type="email" placeholder="Email" onChange={(e) => {setUsername(e.target.value)}} />
        <input className={classes.input} type="password" placeholder="Password" onChange={(e) => {setPassword(e.target.value)}}/>
        <button className={classes.login} type="submit">Login</button>
      </form>
    </div>
  );
}