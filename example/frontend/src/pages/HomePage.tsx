import React from 'react';
import { Link } from "react-router-dom";
import { Container, Row, Col, Button } from "react-bootstrap";

export default class HomePage extends React.Component {
    render() {
        return (
            <Container>
                <Row className="justify-content-md-center">
                    <Col>
                        <h1 className="display-4 font-weight-normal">JWT Auth Proxy Example</h1>
                        <p className="lead font-weight-normal">This web application is meant to demonstrate the usage of the JWT Auth Proxy. You can use it as a template to build your own awesome and secure application.</p>
                        <Button href="https://github.com/virtualzone/jwt-auth-proxy" className="btn btn-primary">GitHub Page</Button>
                        <Link to="/signup.html" className="btn btn-outline-secondary">Sign up</Link>
                        <Link to="/login.html" className="btn btn-outline-secondary">Log in</Link>
                    </Col>
                </Row>
            </Container>
        );
    }
}
