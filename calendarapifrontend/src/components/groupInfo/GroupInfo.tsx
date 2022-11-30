import { useEffect, useState } from "react";
import { groupInfo, Particapant } from "../../../functions/backendFetch";
import styles from "./GroupInfo.module.css";

export default function GroupInfo({ groupId, token, othersCanAdd }: { groupId: string, token: string, othersCanAdd: boolean }) {
    const [data, setData] = useState<any>(null);
    const daytostring = ["N/A", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
    const montostring = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
    useEffect(() => {
        groupInfo(groupId, token).then((data) => {
            if (data.error || !data.data) {
                console.log(data.error);
                setData(<div className={styles.error}>
                    <p>Error: {data.error}</p>
                </div>)
            } else {
                const date = new Date(data.data.created);
                const newDate = `${daytostring[date.getDay()]}, ${montostring[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()}.`;
                setData(
                <div className={styles.groupInfo}>
                    <p className={styles.info}>Info</p>
                    <ul className={styles.infoList}>
                        <li><span>Owner: </span>{data.data.owner}{data.data.yourowner ? " (you)" : null}</li>
                        <li>
                            <span>Particapants: </span>
                            <ul className={styles.users}>
                                {data.data.particapants.map((particapant: Particapant) => <li key={particapant.id}>
                                    <div className={styles.particapantListItem} >
                                        <p>{particapant.name}</p>
                                        {data.data?.yourowner && data.data.yourowner.ownerId !== particapant.id ? <button>Kick Out</button>:null}
                                    </div>
                                </li>)}
                            </ul>
                        </li>
                        {data.data.yourowner && data.data.yourowner.pending_particapants ? <li>
                            <span>Pending Particapants: </span>
                            <ul className={styles.users}>
                                {data.data.yourowner.pending_particapants.map((particapant: Particapant) => <li key={particapant.id}>
                                    <div className={styles.particapantListItem} >
                                        <p>{particapant.name}</p>
                                        {data.data?.yourowner ? <div className={styles.accept}><button>Accept</button><button>Decline</button></div>:null}
                                    </div>
                                </li>)}
                            </ul>
                        </li> : null}
                        <li><span>Date Created: </span>{newDate}</li>
                        <li><span>Group Id: </span>{data.data.group_id}</li>
                        <li><span>Particapants can add jobs: </span>{othersCanAdd ? "Yes": "No"}</li>
                        <li><span>About Group: </span>{data.data.about_group}</li>
                    </ul>
                    <div className={styles.buttons}>
                        {data.data.yourowner ? <button className={styles.leaveGroup}>Delete Group</button> : null}
                        <button className={styles.leaveGroup}>Leave Group</button>
                    </div>
                </div>);
            }
        });
    }, []);

    return data;
}