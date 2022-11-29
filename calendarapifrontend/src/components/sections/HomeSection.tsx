import HomeJob from "../HomeJob";
import styles from "./HomeSection.module.css";

export default function HomeSection({ name }: { name: string }) {
  return (
    <div className={styles.home}>
        <div className={styles.greeting}>
            <p className={styles.welcome}>Welcome</p>
            <p className={styles.name}>{name}</p>
        </div>
        <div className={styles.comingUp}>
            <p className={styles.title}>Coming up</p>
            <div className={styles.comingUpList}>
                <HomeJob name="Babysitting" client="Galina Shapoval" time="Tomorrow 10:30 PM"/>
            </div>
        </div>
    </div>
  )
}