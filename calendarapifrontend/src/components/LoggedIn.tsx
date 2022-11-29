import { useEffect, useState } from "react";
import { fetchGroups } from "../../functions/backendFetch";
import Groups, { GroupsType } from "./Groups";
import PendingGroups, { PendingGroupsType } from "./PendingGroups";
import styles from "./LoggedIn.module.css";
import Navigation, { Locations } from "./Navigation";
import HomeSection from "./sections/HomeSection";
import GroupSection from "./sections/GroupSection";

export default function LoggedIn({ name, email, token, logout }: { name: string, email: string, token: string, logout: () => void }) {
    const [groups, setGroups] = useState<Array<GroupsType>>([]);
    const [pendingGroups, setPendingGroups] = useState<Array<PendingGroupsType>>([]);
    const [selected, setSelected] = useState<Locations>("home");

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
        const doc = document.querySelector("#scroll");
        if (!doc) return;
        doc.scrollLeft = doc.getBoundingClientRect().width * (selected==="addJob"?3:selected==="groups"?2:selected==="calendar"?1:0);
    }, [selected]);
    
  return (
    <div className={styles.main}>
        <div className={styles.sections} id="scroll">
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
            <GroupSection groups={groups} pendingGroups={pendingGroups}/>
            <section>section 4</section>
            <section>section 5</section>
        </div>
        <Navigation selected={selected} navigate={setSelected}/>
    </div>
  );
}
