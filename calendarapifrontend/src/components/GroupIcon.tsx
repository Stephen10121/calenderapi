import styles from "./GroupIcon.module.css";

export default function GroupIcon({ name, owner, id, click, othersCanAdd }: { name: string, owner: string, id: string, click: (id: string, name: string, othersCanAdd: boolean) => void, othersCanAdd: boolean }) {
  return (
    <button className={styles.main} onClick={() => click(id, name, othersCanAdd)}>
        <p className={styles.name}>{name}</p>
        <p className={styles.owner}>{owner}</p>
    </button>
  )
}
