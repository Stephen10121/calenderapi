import styles from "./GroupIcon.module.css";

export default function GroupIcon({ name, owner, id, click }: { name: string, owner: string, id: string, click: (id: string) => void }) {
  return (
    <button className={styles.main} onClick={() => click(id)}>
        <p className={styles.name}>{name}</p>
        <p className={styles.owner}>{owner}</p>
    </button>
  )
}
