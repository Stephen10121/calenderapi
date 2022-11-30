import { useEffect, useState } from "react";
import styles from "./SlideUp.module.css";

type Border = "black" | "red" | "blue";

export interface SlideUpData {
  show: boolean;
  header: string;
  children: any;
  border: Border;
}

export default function SlideUp({ header, children, close, border }: { header: string, children: any, close: () => void, border: Border}) {
  const [top, setTop] = useState("100vh");
  const [borderColor, setBorderColor] = useState("#000000")
  useEffect(() => {
    setBorderColor(border==="blue"?"#3A9FE9":border==="black"?"#000000":border==="red"?"#EE3F3F":"#000000")
    setTimeout(() => {
      setTop("20px");
    }, 1);
  }, []);
  return (
    <div className={styles.main} style={{ top, border: `2px solid ${borderColor}` }}>
        <div className={styles.header}>
            <p className={styles.headerTitle}>{header}</p>
            <button onClick={() => {
              setTop("100vh");
              setTimeout(() => {
                close()
              }, 250);
            }} className={styles.closeButton}>
              <img src="/closecircle.png" alt="Close" />
            </button>
        </div>
        <div className={styles.main2}>
            {children}
        </div>
    </div>
  )
}
