import styles from "./HomeJob.module.css";

export default function HomeJob({ name, client, time }: { name: string, client: string, time: string }) {
  return (
    <div className={styles.box}>
        <div>
            <p className={styles.name}>{name}</p>
            <p className={styles.client}>{client}</p>
        </div>
        <p className={styles.time}>{time}</p>
    </div>
  )
}
