import {
  DocumentScanner,
  Info,
  Security,
  PowerSettingsNew,
} from "@mui/icons-material";
import { Box, Dialog, Grid, Paper, styled } from "@mui/material";
import { Link } from "react-router-dom";
import { appName, theme } from "../constant/config";
import { Component } from "react";
import { Logout } from "../../wailsjs/go/service/Auth";

const AboutDialog = (props: { isOpen: boolean; onClose: VoidFunction }) => {
  return (
    <Dialog
      open={props.isOpen}
      onClose={props.onClose}
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "100%",
        textAlign: "center",
        borderRadius: 20,
      }}
    >
      <div
        style={{
          backgroundColor: theme.palette.primary.main,
          color: "white",
          padding: 20,
        }}
      >
        <div>{appName.toUpperCase()} v1</div>
        <div
          style={{
            padding: "20px 10px",
          }}
        >
          Automated Duties And Excise Payment Verification And Information
          System
        </div>
      </div>
    </Dialog>
  );
};

const MyButton = styled(Paper)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  padding: theme.spacing(4),
  textAlign: "center",
  borderRadius: 4,
  color: "white",
  cursor: "pointer",
  textDecoration: "none",
}));

class HomeScreen extends Component {
  state = {
    showAbout: false,
  };

  handleNotAvaliable() {
    const n = document.getElementById("notify")!;
    n.innerHTML = "Not avaliable at this time";
    setTimeout(() => {
      n.innerHTML = "";
    }, 2000);
  }

  render() {
    return (
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          height: "100%",
        }}
      >
        <Box
          sx={{
            borderRadius: 5,
            padding: 2,
            color: "white",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(2, 1fr)",
              columnGap: 20,
              rowGap: 20,
            }}
          >
            <Link to={"/scan"} style={{ textDecoration: "none" }}>
              <MyButton>
                <DocumentScanner sx={{ fontSize: 50 }} />
                <div>SCAN</div>
              </MyButton>
            </Link>

            <MyButton onClick={() => this.handleNotAvaliable()}>
              <Security sx={{ fontSize: 50 }} />
              <div>SECURITY</div>
            </MyButton>

            <MyButton onClick={() => this.setState({ showAbout: true })}>
              <Info sx={{ fontSize: 50 }} />
              <div>ABOUT</div>
            </MyButton>

            <MyButton onClick={Logout}>
              <PowerSettingsNew sx={{ fontSize: 50 }} />
              <div>LOGOUT</div>
            </MyButton>
          </div>
        </Box>
        <div
          id="notify"
          style={{
            transition: "ease-in .5s",
            color: "white",
          }}
        ></div>
        <AboutDialog
          isOpen={this.state.showAbout}
          onClose={() =>
            this.setState({
              showAbout: false,
            })
          }
        />
      </div>
    );
  }
}

export default HomeScreen;
