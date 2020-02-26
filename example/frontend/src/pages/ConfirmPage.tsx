import React from 'react';
import { Container, Row, Col, Button } from "react-bootstrap";
import Ajax from '../Ajax';

interface ConfirmPageState {
    isLoading: boolean
    hasSuccess: boolean
    hasDataError: boolean
}

export default class ConfirmPage extends React.Component<{}, ConfirmPageState> {
    id: string | null;

    constructor(props: any) {
        super(props);
        let query = new URLSearchParams(window.location.search);
        this.id = query.get("id");
        this.state = {
            isLoading: true,
            hasSuccess: false,
            hasDataError: false
        };
    }

    componentDidMount() {
        this.performConfirm();
    }

    performConfirm() {
        Ajax.postData("/auth/confirm/" + this.id).then(res => {
            if (res.status === 204) {
                this.setState({
                    isLoading: false,
                    hasSuccess: true,
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
                        <p className="lead font-weight-normal">We are updating your data...</p>
                    </Col> : null }
                    { this.state.hasDataError ? <Col lg="7">
                        <h1 className="display-4 font-weight-normal">Something went wrong.</h1>
                        <p className="lead font-weight-normal">Your link may either be outdated or the action has already been performed.</p>
                        <Button href="/" className="btn btn-primary">Return to the home page</Button>
                    </Col> : null }
                    { this.state.hasSuccess ? <Col lg="7">
                        <h1 className="display-4 font-weight-normal">That worked!</h1>
                        <Button href="/" className="btn btn-primary">Return to the home page</Button>
                    </Col> : null }
                </Row>
            </Container>
        );
    }
}
