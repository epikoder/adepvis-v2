import { Box, Button, IconButton, TextField, Typography } from "@mui/material";
import { ChangeEvent, Component } from "react";
import { useNavigate } from "react-router-dom";
import { appName } from "../constant/config";
import Logo from "../logo.jpg";
import { VisibilityOff, Visibility } from "@mui/icons-material";
import { Login } from "../../wailsjs/go/service/Auth";
import { LogDebug, LogInfo } from "../../wailsjs/runtime/runtime";
import Copyright from "../components/copyright";
import { invoke } from "../helper";

interface Props {}
interface State {
  svn: string;
  username: string;
  password: string;
  isPasswordHidden: boolean;
  loading: boolean;
  message: string[];
}

class LoginScreen extends Component<Props, State> {
  state: State = {
    svn: "88716",
    username: "test.officer",
    password: "KqIBmSUEzxFt",
    isPasswordHidden: true,
    loading: false,
    message: [],
  };

  handleChange = (current: ChangeEvent<HTMLInputElement>) => {
    this.setState({
      ...this.state,
      [current.target.name]: current.target.value,
    });
  };

  handleSubmit = async () => {
    this.setState({ message: [], loading: true });
    const res = await invoke(() =>
      Login({
        svn: this.state.svn,
        username: this.state.username,
        password: this.state.password,
      })
    );

    if (res.Value === true) return this.setState({ loading: false });
    this.setState({
      message: [
        ...(res.Err?.message === "Unauthorized"
          ? [res.Err.message]
          : ["Error occured"]),
      ],
      loading: false,
    });
  };

  render() {
    return (
      <div
        style={{
          height: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Box
          sx={{
            bgcolor: "primary.main",
            color: "white",
            display: "flex",
            flexDirection: "row",
            borderRadius: 4,
            height: "100%",
            width: "100%",
          }}
        >
          <div
            style={{
              width: "60%",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <div
              style={{
                position: "relative",
                left: "50%",
              }}
            >
              <img src={Logo} width="200px" className="logo" alt={appName} />
            </div>
            <Typography
              fontSize={18}
              color="white"
              sx={{
                position: "absolute",
              }}
            >
              {appName.toUpperCase()}
            </Typography>
          </div>
          <div
            style={{
              width: "40%",
              backgroundColor: "white",
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              borderRadius: "0px 16px 16px 0px",
            }}
          >
            <div
              style={{
                padding: 2,
                borderRadius: 10,
                marginLeft: 20,
              }}
            >
              <div
                style={{
                  margin: "20px 0px",
                  padding: 0.5,
                  backgroundColor: "red",
                  borderRadius: 5,
                  visibility:
                    this.state.message.length > 0 ? "visible" : "hidden",
                }}
              >
                {this.state.message.map((m) => (
                  <Typography key={m}>{m}</Typography>
                ))}
              </div>
              <div>
                <Typography color="#0b0b22">OFFICER LOGIN</Typography>
              </div>
              <div
                style={{
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <TextField
                  size="small"
                  label="Service Number"
                  placeholder="11111"
                  value={this.state.svn}
                  onChange={this.handleChange}
                  type={"number"}
                  name="svn"
                  sx={{
                    margin: 0.5,
                  }}
                />
                <TextField
                  size="small"
                  label="Username"
                  value={this.state.username}
                  type={"text"}
                  name="username"
                  onChange={this.handleChange}
                  sx={{
                    margin: 0.5,
                  }}
                />
                <div
                  style={{
                    position: "relative",
                  }}
                >
                  <TextField
                    size="small"
                    label="Password"
                    value={this.state.password}
                    type={this.state.isPasswordHidden ? "password" : "text"}
                    name="password"
                    onChange={this.handleChange}
                    sx={{
                      margin: 0.5,
                    }}
                  />
                  <IconButton
                    sx={{
                      position: "absolute",
                      right: "5px",
                      top: "3px",
                    }}
                    onClick={() =>
                      this.setState({
                        isPasswordHidden: !this.state.isPasswordHidden,
                      })
                    }
                  >
                    {this.state.isPasswordHidden ? (
                      <Visibility />
                    ) : (
                      <VisibilityOff />
                    )}
                  </IconButton>
                </div>
                <div>
                  <Button
                    disabled={this.state.loading}
                    size="small"
                    variant="outlined"
                    onClick={this.handleSubmit}
                    sx={{ margin: 0.5 }}
                  >
                    SUBMIT
                  </Button>
                </div>
                <div
                  style={{
                    position: "absolute",
                    bottom: 10,
                    justifyContent: "center",
                    display: "flex",
                  }}
                >
                  <Copyright color={"black"} />
                </div>
              </div>
            </div>
          </div>
        </Box>
      </div>
    );
  }
}

export default LoginScreen;
