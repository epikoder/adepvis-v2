import { ArrowBackIos, Replay, WifiTethering } from "@mui/icons-material";
import { useNavigate } from "react-router-dom";
import { appName, theme } from "../constant/config";
import Logo from "../logo.jpg";
import { Box, Button, Dialog } from "@mui/material";
import QrCodeScanner from "../components/qrcode_scanner";
import {
  SetId,
  StartScanner,
  StopScanner,
} from "../../wailsjs/go/service/Scanner";
import { EventsOn, LogError, LogInfo } from "../../wailsjs/runtime/runtime";
import { useEffect, useState } from "react";
import { service } from "../../wailsjs/go/models";

const TR = (props: { name: string; value: any }) => {
  return (
    <tr>
      <td
        style={{
          textAlign: "start",
          fontWeight: "bold",
        }}
      >
        {props.name}
      </td>
      <td
        style={{
          textAlign: "end",
        }}
      >
        {props.value}
      </td>
    </tr>
  );
};

const formatCurrrency = function (
  value: number,
  n: number,
  x: number,
  s: string,
  c?: string
) {
  var re = "\\d(?=(\\d{" + (x || 3) + "})+" + (n > 0 ? "\\D" : "$") + ")",
    num = value.toFixed(Math.max(0, ~~n));
  return (c ? num.replace(".", c) : num).replace(
    new RegExp(re, "g"),
    "$&" + (s || ",")
  );
};

const VehicleInfo = (props: { vehicle: any; client: any }) => {
  const vehicle = props.vehicle;
  const client = props.client;

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-around",
        width: "80%",
        height: "80%",
      }}
    >
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
        }}
      >
        <div
          style={{
            border: "2px solid black",
          }}
        >
          <img width={100} height={100} alt="Avatar" src={Logo} />
        </div>
        <Box
          sx={{
            bgcolor: "primary.main",
            borderRadius: 4,
            color: "white",
            padding: 1.5,
            width: "50%",
          }}
        >
          <table style={{ width: "100%" }}>
            <tbody>
              <TR name="Name" value={client.name} />
              <TR name="Tax identification no." value={client.tin} />
              <TR name="Address" value={client.address} />
            </tbody>
          </table>
        </Box>
      </div>
      <Box
        sx={{
          padding: 1.5,
          bgcolor: "primary.main",
          borderRadius: 4,
          color: "white",
        }}
      >
        <table
          style={{
            width: "100%",
          }}
        >
          <tbody>
            <TR name="Chasis Number" value={vehicle.chasis_number} />
            <TR name="Vehicle Identification Number" value={vehicle.vin} />
            <TR name="Type" value={vehicle.type} />
            <TR name="Model" value={vehicle.model} />
            <TR name="Year" value={vehicle.year} />
            <TR name="Purpose" value={vehicle.purpose} />
            <TR
              name="Duty Fee"
              value={`₦${formatCurrrency(
                ((vehicle.cost * vehicle.duty_fee) / 100) *
                  vehicle.current_rate,
                2,
                3,
                ","
              )}`}
            />
          </tbody>
        </table>
        <div
          style={{
            display: "flex",
            flexDirection: "row",
            justifyContent: "space-around",
          }}
        >
          <div
            style={{
              padding: 10,
              borderRadius: 10,
              backgroundColor: vehicle.stolen ? "red" : "green",
              color: "white",
            }}
          >
            {vehicle.stolen ? "STOLEN" : "CLEAN"}
          </div>
        </div>
      </Box>
    </div>
  );
};

const ScanScreen = () => {
  const [scanStatus, setScanStatus] = useState({
    message: "",
    status: false,
  });
  const [scanComplete, setScanComplete] = useState(false);
  const [vehicleInfo, setVehicleinfo] = useState({});
  const [clientInfo, setClientInfo] = useState({});
  const [scanOption, setScanOption] = useState<string | undefined>(undefined);
  const [runing, setRuning] = useState(false);
  const navigate = useNavigate();

  const startScanner = (value: "camera" | "card") => {
    setScanComplete(false);
    LogInfo(value);
    StartScanner({
      Camera: value == "camera",
      Card: value === "card",
    })
      .then((res) => {
        if (res === undefined) return;
        let data = JSON.parse(res);
        console.log(data);
        if (data.status === "failed")
          return setScanStatus({
            message: data.message,
            status: false,
          });
        setScanComplete(true);
        setScanStatus({
          message: "",
          status: true,
        });
        setVehicleinfo(data.data.vehicle);
        setClientInfo(data.data.user.client);
      })
      .catch((err) => {
        LogError(err);
        if (err === undefined) return;
        setScanStatus({
          message: err,
          status: false,
        });
      });
    return setRuning(false);
  };

  useEffect(() => {
    EventsOn("scanner", (state) => {
      LogInfo(`State: ${runing}`);
      if (runing) return;
      setScanStatus({
        message: state.Message,
        status: state.Status,
      });
    });
  }, []);

  const onSelectOption = (value: "card" | "camera") => {
    setScanOption(value);
    startScanner(value);
    return setRuning(true);
  };

  const back = async () => {
    if (runing) {
      await StopScanner();
      return setRuning(false);
    }
    return navigate("/");
  };

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
      <div
        onClick={() => back()}
        style={{
          position: "absolute",
          top: 10,
          left: 10,
          padding: 10,
          borderRadius: 10,
          display: "flex",
          justifyContent: "space-around",
          backgroundColor: theme.palette.primary.main,
          color: "white",
          cursor: "pointer",
        }}
      >
        {!runing ? <ArrowBackIos /> : null}
        <div>{runing ? "STOP" : "BACK"}</div>
      </div>
      {!scanComplete && scanOption === "card" && (
        <div
          style={{
            position: "relative",
            left: -50,
            top: -100,
          }}
        >
          <div
            style={{
              position: "absolute",
              top: -90,
            }}
          >
            <WifiTethering
              sx={{
                color: "#9b9b9b",
                fontSize: 150,
              }}
            />
          </div>
          <div
            style={{
              backgroundColor: "#ededed",
              width: 120,
              height: 200,
              left: 15,
              top: -20,
              borderRadius: 10,
              position: "absolute",
              display: "flex",
              justifyContent: "center",
              flexDirection: "column",
              alignItems: "center",
              color: "black",
            }}
          >
            {appName.toUpperCase()}
          </div>

          <div
            style={{
              position: "absolute",
              top: 190,
              textAlign: "center",
              color: scanStatus.status ? "green" : "red",
              width: 150,
            }}
          >
            {scanStatus.message}
          </div>
          {!scanStatus.status && (
            <div
              style={{
                position: "absolute",
                top: 250,
                textAlign: "center",
                width: 150,
                cursor: "pointer",
              }}
              onClick={() => setScanOption(undefined)}
            >
              <Replay sx={{ color: "white" }} />
            </div>
          )}
        </div>
      )}
      {scanComplete && (
        <VehicleInfo vehicle={vehicleInfo} client={clientInfo} />
      )}
      {!scanComplete && scanOption === "camera" && (
        <div>
          <QrCodeScanner onCompleted={async (text) => await SetId(text)} />
        </div>
      )}
      {!scanComplete && scanOption === "camera" && (
        <div>
          <div
            style={{
              textAlign: "center",
              color: scanStatus.status ? "green" : "red",
            }}
          >
            {scanStatus.message}
          </div>
          {!scanStatus.status && scanStatus.message !== "" && (
            <div
              style={{
                cursor: "pointer",
              }}
              onClick={async () => {
                await StopScanner();
                setScanOption(undefined);
              }}
            >
              <Replay sx={{ color: "white" }} />
            </div>
          )}
        </div>
      )}
      {scanComplete && (
        <Button
          sx={{
            color: "white",
          }}
          onClick={() => setScanOption(undefined)}
          variant="outlined"
        >
          SCAN AGAIN
        </Button>
      )}
      <Dialog
        open={scanOption === undefined}
        sx={{
          padding: 4,
        }}
        onClose={() => setScanOption("")}
      >
        <Button
          onClick={() => onSelectOption("card")}
          sx={{
            height: "70px",
            width: "200px",
            ":hover": {
              backgroundColor: "#656583",
              color: "white",
            },
            ":focus": {
              backgroundColor: "#656583",
              color: "white",
            },
          }}
        >
          CARD
        </Button>
        <Button
          onClick={() => onSelectOption("camera")}
          sx={{
            height: "70px",
            width: "200px",
            ":hover": {
              backgroundColor: "#656583",
              color: "white",
            },
            ":focus": {
              backgroundColor: "#656583",
              color: "white",
            },
          }}
        >
          CAMERA
        </Button>
      </Dialog>
    </div>
  );
};

export default ScanScreen;
