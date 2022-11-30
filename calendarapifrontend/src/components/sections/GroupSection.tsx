import { useState } from "react";
import GroupIcon from "../GroupIcon";
import GroupInfo from "../groupInfo/GroupInfo";
import { GroupsType } from "../Groups";
import { PendingGroupsType } from "../PendingGroups";
import SlideUp, { SlideUpData } from "../slideUp/SlideUp";
import styles from "./GroupSection.module.css";

export default function GroupSection({ groups, pendingGroups, token }: { groups: GroupsType[], pendingGroups: PendingGroupsType[], token: string }) {
    const [showSlideUp, setShowSlideUp] = useState<SlideUpData>({show: false, header: "N/A", children: null, border:"black"});

    function groupClicked(groupId: string, name: string, othersCanAdd: boolean) {
        setShowSlideUp({ show: true, header: name, children: <GroupInfo token={token} groupId={groupId} othersCanAdd={othersCanAdd}/>, border: "black" });
    }

    return (
        <div className={styles.home}>
            {showSlideUp.show ? <SlideUp border={showSlideUp.border} close={() => setShowSlideUp({...showSlideUp, show: false})} header={showSlideUp.header}>{showSlideUp.children}</SlideUp> : null}
            <div className={styles.greeting}>
                <p className={styles.welcome}>Groups</p>
            </div>
            <div className={styles.comingUp}>
                <p className={styles.title}>Joined</p>
                <div className={styles.comingUpList}>
                    {groups ? groups.map((group) => <GroupIcon key={group.groupId} id={group.groupId} name={group.groupName} owner={group.groupOwner} othersCanAdd={group.othersCanAdd} click={groupClicked}/>) : <div className={styles.nogroup}><p>No Groups</p></div>}
                </div>
            </div>
            {pendingGroups ? 
            <div className={styles.available}>
                <p className={styles.title}>Pending</p>
                <div className={styles.comingUpList}>
                {pendingGroups.map((group) => <GroupIcon key={group.groupId} id={group.groupId} name={group.groupName} othersCanAdd={false} owner="Anonymous" click={console.log}/>)}
                </div>
            </div> : null }
        </div>
    )
}