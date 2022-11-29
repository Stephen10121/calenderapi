import { useEffect, useState } from "react";
import { fetchGroups } from "../../functions/backendFetch";
import Groups, { GroupsType } from "./Groups";
import PendingGroups, { PendingGroupsType } from "./PendingGroups";
import styles from "./LoggedIn.module.css";
import Navigation, { Locations } from "./Navigation";
import HomeSection from "./sections/HomeSection";

export default function LoggedIn({ name, email, token, logout }: { name: string, email: string, token: string, logout: () => void }) {
    const [groups, setGroups] = useState<Array<GroupsType>>([]);
    const [pendingGroups, setPendingGroups] = useState<Array<PendingGroupsType>>([]);
    const [selected, setSelected] = useState<Locations>("home");
    const [width, setWidth] = useState(100);
    const [marginLeftSet, setMarginLeftSet] = useState(0);

    useEffect(() => {
        fetchGroups(token).then((data) => {
            if (data.error || !data.data) {
                console.log(data.error);
                return
            }
            setGroups(data.data.groups);
            setPendingGroups(data.data.pendingGroups);
          });
    }, []);

    useEffect(() => {
        if (selected === "home") {
            setWidth(100);
            setMarginLeftSet(0);
        }
        else if (selected === "calendar") {
            setWidth(200);
            setMarginLeftSet(100);
        }
        else if (selected === "groups") {
            setWidth(300);
            setMarginLeftSet(200);
        }
        else if (selected === "addJob") {
            setWidth(400);
            setMarginLeftSet(300);
        }
    }, [selected]);

  return (
    <div className={styles.main}>
        <div className={styles.sections} style={{ width: `${width}vw`, marginLeft: `-${marginLeftSet}vw` }}>
            <HomeSection name={name}/>
            <section>
                <div className="info">
                    <h1>Hello {name}</h1>
                    <h2>Contacts: {email}</h2>
                </div>
                <div className="group-section">
                    <Groups groups={groups}/>
                    <PendingGroups groups={pendingGroups}/>
                </div>
                {selected}
                <button className="logout" onClick={logout}>Logout</button>
            </section>
            <section>section 3</section>
            <section>section 4</section>
            <section>section 5</section>
        </div>
        <Navigation selected={selected} navigate={setSelected}/>
    </div>
  );
}
