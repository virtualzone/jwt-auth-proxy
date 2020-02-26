import React from 'react';
import './App.css';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import { withRouter } from "react-router";
import { Navbar, Nav } from "react-bootstrap";
import HomePage from './pages/HomePage'
import SignupForm from './pages/SignupForm'
import LoginForm from './pages/LoginForm'
import ConfirmPage from './pages/ConfirmPage'
import DashboardPage from './pages/DashboardPage';
import Ajax from './Ajax';

function RouteWithSubRoutes(route) {
  return (
    <Route
      path={route.path}
      render={props => (
        <route.component {...props} routes={route.routes} />
      )}
    />
  );
}

const NavHeader = props => {
  function handleLogout(event) {
    event.preventDefault();
    let data = {
      "refreshToken": window.sessionStorage.getItem("refreshToken")
    };
    Ajax.postData("/auth/logout", data).then(res => {
      window.sessionStorage.removeItem("refreshToken");
      window.sessionStorage.removeItem("accessToken");
      window.location.href = "/";
      return;
    }).catch((e) => {});
  }

  const { location } = props;
  return (
    <Navbar bg="dark" variant="dark">
      <Nav activeKey={location.pathname}>
        <Nav.Link href="/">Home</Nav.Link>
        {window.sessionStorage.getItem("accessToken") == null ? <Nav.Link href="/signup.html">Sign Up</Nav.Link> : null}
        {window.sessionStorage.getItem("accessToken") == null ? <Nav.Link href="/login.html">Log In</Nav.Link> : null}
        {window.sessionStorage.getItem("accessToken") != null ? <Nav.Link href="/dashboard.html">Dashboard</Nav.Link> : null}
        {window.sessionStorage.getItem("accessToken") != null ? <Nav.Link onClick={handleLogout}>Log Out</Nav.Link> : null}
      </Nav>
    </Navbar>
  );
};
const NavHeaderWithRouter = withRouter(NavHeader);
const routes = [
  {
    path: "/",
    exact: true,
    component: HomePage
  },
  {
    path: "/signup.html",
    component: SignupForm
  },
  {
    path: "/login.html",
    component: LoginForm
  },
  {
    path: "/confirm.html",
    component: ConfirmPage
  },
  {
    path: "/dashboard.html",
    component: DashboardPage
  }
];

function checkRenewToken() {
  if ((window.sessionStorage.getItem("accessToken") != null) && (window.sessionStorage.getItem("refreshToken") != null)) {
    Ajax.postData("/auth/refresh", {"refreshToken": window.sessionStorage.getItem("refreshToken")}).then(res => {
      if (res.status === 200) {
        if (res.json.accessToken) {
          window.sessionStorage.setItem("accessToken", res.json.accessToken);
        }
        if (res.json.refreshToken) {
          window.sessionStorage.setItem("refreshToken", res.json.refreshToken);
        }
        console.log("Access Token refreshed.");
      }
    });
  }
}

window.setInterval(checkRenewToken, 60*1000);

function App() {
  return (
    <div className="App">
      <Router>
        <NavHeaderWithRouter />
        <Switch>
          {routes.map((route, i) => (
            <RouteWithSubRoutes key={i} {...route} />
          ))}
        </Switch>
      </Router>
    </div>
  );
}

export default App;
