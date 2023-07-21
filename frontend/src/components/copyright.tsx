import { GitHub } from "@mui/icons-material";

export default function Copyright({ color = "white" }: { color?: string }) {
  return (
    <div style={{ fontSize: "12px", alignItems: "center", display: "flex" }}>
      <GitHub htmlColor={color} />
      <span style={{ paddingLeft: "8px", color: color }}>Epikoder</span>
    </div>
  );
}
