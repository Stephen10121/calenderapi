import { useEffect, useState } from "react";
import styles from "./SlideUp.module.css";

export default function SlideUp({ header, children, close }: { header: string, children: any, close: () => void }) {
  const [top, setTop] = useState("100vh");
  useEffect(() => {
    setTimeout(() => {
      setTop("20px");
    }, 1);
  }, []);
  return (
    <div className={styles.main} style={{ top }}>
        <div className={styles.header}>
            <p>{header}</p>
            <button onClick={() => {
              setTop("100vh");
              setTimeout(() => {
                close()
              }, 250);
            }}>
              <img src="/closecircle.png" alt="Close" />
            </button>
        </div>
        <div>
            {children}
        </div>
    </div>
  )
}
