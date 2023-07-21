import { GitHub } from "@mui/icons-material";
import Spinner from "../components/spinner";
import { Box, Typography } from "@mui/material";
import Copyright from "../components/copyright";

const LaunchScreen = () => {
  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Typography
        sx={{
          color: "white",
          fontSize: "32px",
        }}
      >
        <span style={{ fontWeight: "bolder" }}>A</span>depvis
      </Typography>
      <div
        style={{
          position: "absolute",
          bottom: 30,
        }}
      >
        <Spinner />
        <Copyright />
      </div>
    </Box>
  );
};

export default LaunchScreen;
