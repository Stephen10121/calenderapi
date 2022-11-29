import styles from "./Navigation.module.css";

export type Locations = "home" | "calendar" | "groups" | "addJob"

export default function Navigation({ selected, navigate } : { selected: Locations, navigate: (arg0: any) => any }) {
  return (
    <div className={styles.navigation}>
        <button className={selected === "home" ? styles.selected : ""} onClick={() => {navigate("home")}}>
            <img src="/home.png" alt="Home" />
        </button>
        <button className={selected === "calendar" ? styles.selected : ""} onClick={() => {navigate("calendar")}}>
            <img src="/calendar.png" alt="Calendar" />
        </button>
        <button className={selected === "groups" ? styles.selected : ""} onClick={() => {navigate("groups")}}>
            <img src="/groups.png" alt="Groups" />
        </button>
        <button className={selected === "addJob" ? styles.selected : ""} onClick={() => {navigate("addJob")}}>
            <img src="/addjob.png" alt="Add Job" />
        </button>
        <button>
            <img src="/avatar.png" alt="My Account" />
        </button>
    </div>
  )
}