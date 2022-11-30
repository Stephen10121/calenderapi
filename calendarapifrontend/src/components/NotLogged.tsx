import { FormEvent, useState } from "react";
import { login } from "../../functions/backendFetch";
import classes from "./NotLogged.module.css";

interface ResponseData {
  email: string
  name: string
  token: string
}

export default function NotLogged({ loggedIn }: { loggedIn: (email: string, password: string, token: string) => void}) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  // const responseGoogle = (response: any) => {
  //   console.log(response);
  // }
  // const clientId = "665789630027-g9gogfmb2qo3s4egs0boghfhc2so267i.apps.googleusercontent.com";


  // useEffect(() => {
  //   // global google
  //   //@ts-ignore
  //   google.accounts!.id.initialize({
  //     clientId,
  //     callback: responseGoogle
  //   });
  //   const googleDiv = document.querySelector("#googleSignup");
  //   //@ts-ignore
  //   google.accounts.id.renderButton(googleDiv, {
  //     theme: "outline",
  //     size: "large",
  //     type: "standard"
  //   });

  //   //@ts-ignore
  //   google.accounts.id.prompt();
  // }, []);

  function formSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    login(username, password).then((data) => {
      if (data.error) {
        setLoading(false);
        setError(data.error);
        return
      }
      const response = data.data as ResponseData;
      document.cookie = `G_CAL=${response.token}`;
      setLoading(false);
      loggedIn(response.email, response.name, response.token);
    });
  }

  return (
    <div className={classes.notLogged}>
      <h1 className={classes.welcome}>Welcome</h1>
      {error.length !== 0 ? <p className={classes.error}>{error}</p> : null}
      <form onSubmit={formSubmit}>
        <input className={classes.input} type="email" placeholder="Email" onChange={(e) => {setUsername(e.target.value)}} />
        <input className={classes.input} type="password" placeholder="Password" onChange={(e) => {setPassword(e.target.value)}}/>
        <button className={classes.login} type="submit">{loading ? <div className="preloader"><span className="loader2"></span></div>: "Login"}</button>
        <button className={classes.input2}><img src="/google.png" alt="Google Icon" /><p>Continue With Google</p></button>
      </form>
      <div id="googleSignup"></div>,
    </div>
  );
}