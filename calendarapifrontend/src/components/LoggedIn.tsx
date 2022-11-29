import { useEffect, useState } from "react";
import { fetchGroups } from "../../functions/backendFetch";
import Groups, { GroupsType } from "./Groups";
import PendingGroups, { PendingGroupsType } from "./PendingGroups";

export default function LoggedIn({ name, email, token, logout }: { name: string, email: string, token: string, logout: () => void }) {
    const [groups, setGroups] = useState<Array<GroupsType>>([]);
    const [pendingGroups, setPendingGroups] = useState<Array<PendingGroupsType>>([]);

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

  return (
    <div>
        <div className="info">
            <h1>Hello {name}</h1>
            <h2>Contacts: {email}</h2>
        </div>
        <div className="group-section">
            <Groups groups={groups}/>
            <PendingGroups groups={pendingGroups}/>
        </div>
        <button className="logout" onClick={logout}>Logout</button>
    </div>
  );
}
