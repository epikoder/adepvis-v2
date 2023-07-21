import { createTheme } from "@mui/material";

const appName = "Adepvis";
const apiURL = "https://adepvis.herokuapp.com";

const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#171622",
    },
    secondary: {
      main: "#212130",
    },
  },
  typography: {
    fontFamily: '"Roboto", "Oswald"',
    fontSize: 12,
  },
});
export { apiURL, theme, appName };
