import { useState } from "react";
import GroupIcon from "../GroupIcon";
import { GroupsType } from "../Groups";
import { PendingGroupsType } from "../PendingGroups";
import SlideUp from "../slideUp/SlideUp";
import styles from "./GroupSection.module.css";

export default function GroupSection({ groups, pendingGroups }: { groups: GroupsType[], pendingGroups: PendingGroupsType[] }) {
    const [showSlideUp, setShowSlideUp] = useState(false);
    function groupClicked(groupId: string) {
        console.log(groupId);
        setShowSlideUp(true);
    }
    return (
        <div className={styles.home}>
            {showSlideUp ? <SlideUp close={() => setShowSlideUp(false)} header="Test">wow2</SlideUp> : null}
            <div className={styles.greeting}>
                <p className={styles.welcome}>Groups</p>
            </div>
            <div className={styles.comingUp}>
                <p className={styles.title}>Joined</p>
                <div className={styles.comingUpList}>
                    {groups ? groups.map((group) => <GroupIcon key={group.groupId} id={group.groupId} name={group.groupName} owner={group.groupOwner} click={groupClicked}/>) : "No Groups"}
                </div>
            </div>
            {pendingGroups ? 
            <div className={styles.available}>
                <p className={styles.title}>Pending</p>
                <div className={styles.comingUpList}>
                {pendingGroups.map((group) => <GroupIcon key={group.groupId} id={group.groupId} name={group.groupName} owner="Anonymous" click={console.log}/>)}
                </div>
            </div> : null }
        </div>
    )
}