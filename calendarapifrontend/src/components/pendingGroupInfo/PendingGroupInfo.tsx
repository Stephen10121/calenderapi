import styles from "../groupInfo/GroupInfo.module.css";

export default function PendingGroupInfo({ groupId, token }: { groupId: string, token: string }) {
    return(
    <div className={styles.groupInfo}>
        <p className={styles.info}>Info</p>
        <ul className={styles.infoList}>
            <li><span>Group Id: </span>{groupId}</li>
        </ul>
        <div className={styles.buttons}>
            <button className={styles.leaveGroup}>Cancel Request</button>
        </div>
    </div>);
}