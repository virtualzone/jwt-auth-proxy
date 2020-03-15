import React from 'react';
import { Container, Row, Col, Form, Button, Alert } from "react-bootstrap";
import Ajax from '../Ajax';

interface TotpEnableFormState {
    isLoading: boolean
    hasDataError: boolean
    stateInit: boolean
    stateConfirm: boolean
    stateFinish: boolean
    secret: string
    qrCode: string
    otp: string
}

export default class TotpEnableForm extends React.Component<{}, TotpEnableFormState> {
    constructor(props: any) {
        super(props);
        this.state = {
            isLoading: false,
            hasDataError: false,
            stateInit: true,
            stateConfirm: false,
            stateFinish: false,
            secret: "",
            qrCode: "",
            otp: ""
        };
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmitInit = this.handleSubmitInit.bind(this);
        this.handleSubmitConfirm = this.handleSubmitConfirm.bind(this);
    }

    handleChange(event: React.FormEvent) {
        let target = event.target as HTMLInputElement;
        this.setState({[target.name]: target.value} as React.ComponentState );
        event.preventDefault();
    }

    handleSubmitInit(event: React.FormEvent) {
        event.preventDefault();
        this.setState({
            isLoading: true,
            hasDataError: false
        });
        Ajax.postData("/auth/otp/init").then(res => {
            if (res.status === 200) {
                this.setState({
                    isLoading: false,
                    hasDataError: false,
                    stateInit: false,
                    stateConfirm: true,
                    secret: res.json.secret,
                    qrCode: "data:image/png;base64," + res.json.image
                });
                return;
            }
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        }).catch(() => {
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        });
    }

    handleSubmitConfirm(event: React.FormEvent) {
        event.preventDefault();
        this.setState({
            isLoading: true,
            hasDataError: false
        });
        let payload = {
            passcode: this.state.otp
        };
        Ajax.postData("/auth/otp/confirm", payload).then(res => {
            if (res.status === 204) {
                this.setState({
                    isLoading: false,
                    hasDataError: false,
                    stateConfirm: false,
                    stateFinish: true
                });
                return;
            }
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        }).catch(() => {
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        });
    }

    render() {
        return (
            <Container>
                <Row className="justify-content-md-center">
                    <Col lg="5">
                        <h1 className="display-4 font-weight-normal">Security</h1>
                        <Alert variant="danger" show={this.state.hasDataError}>Please verify the data you have entered.</Alert>
                        <Form onSubmit={this.handleSubmitInit} hidden={!this.state.stateInit}>
                            <p>Two-Factor Authentication is not enabled yet. Enable it to add an additional layer of security to your account.</p>
                            <Button variant="primary" type="submit" disabled={this.state.isLoading}>{this.state.isLoading ? 'Loading...' : 'Enable 2FA'}</Button>
                        </Form>
                        <Form onSubmit={this.handleSubmitConfirm} hidden={!this.state.stateConfirm}>
                            <p>Scan the barcode below with your authenticator app, orenter the code below if you can't use the barcode.</p>
                            <p><img src={this.state.qrCode} alt="" /></p>
                            <p>{this.state.secret}</p>
                            <Form.Group controlId="otp">
                                <Form.Label>Six-Digit Code from authenticator app</Form.Label>
                                <Form.Control type="text" pattern="\d*" name="otp" placeholder="Time-based One-time Password (TOTP)" value={this.state.otp} onChange={this.handleChange} minLength={6} maxLength={6} required={true} />
                            </Form.Group>
                            <Button variant="primary" type="submit" disabled={this.state.isLoading}>{this.state.isLoading ? 'Loading...' : 'Finish 2FA Setup'}</Button>
                        </Form>
                        <div hidden={!this.state.stateFinish}>
                            <p>2FA is now enabled for your account.</p>
                        </div>
                    </Col>
                </Row>
            </Container>
        );
    }
}
