import React from 'react';
import { Container, Row, Col } from "react-bootstrap";
import Ajax from '../Ajax';

interface DashboardPageState {
    isLoading: boolean
    hasSuccess: boolean
    hasDataError: boolean
    data: any
}

export default class DashboardPage extends React.Component<{}, DashboardPageState> {
    constructor(props: any) {
        super(props);
        this.state = {
            isLoading: true,
            hasSuccess: false,
            hasDataError: false,
            data: {}
        };
    }

    componentDidMount() {
        this.loadUserData();
    }

    loadUserData() {
        Ajax.get("/api/userinfo").then(res => {
            if (res.status === 200) {
                this.setState({
                    isLoading: false,
                    hasSuccess: true,
                    data: res.json
                });
                return;
            }
            this.setState({
                isLoading: false,
                hasDataError: true,
            });
        }).catch(() => {
            this.setState({
                isLoading: false,
                hasDataError: true,
            });
        });
    }

    render() {
        return (
            <Container>
                <Row className="justify-content-md-center">
                    { this.state.isLoading ? <Col lg="7">
                        <h1 className="display-4 font-weight-normal">Please wait...</h1>
                        <p className="lead font-weight-normal">We are loading your data...</p>
                    </Col> : null }
                    { this.state.hasDataError ? <Col lg="7">
                        <h1 className="display-4 font-weight-normal">Something went wrong.</h1>
                        <p className="lead font-weight-normal">That should not have happened. Please try again.</p>
                    </Col> : null }
                    { this.state.hasSuccess ? <Col lg="7">
                        <h1 className="display-4 font-weight-normal">Welcome!</h1>
                        <p className="lead font-weight-normal">You are logged in as {this.state.data.email}.</p>
                        <p className="lead font-weight-normal">You signed up at {this.state.data.createDate}.</p>
                        <p className="lead font-weight-normal">If you can read the information above, the example is working.</p>
                    </Col> : null }
                </Row>
            </Container>
        );
    }
}