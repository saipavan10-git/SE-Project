import React from 'react';
import { Button } from '@mui/material';
import { Link } from 'react-router-dom';

function Homepage() {
    return (
        <>
            <div>
                <div className='centered'>
                    <h1>COVID BOOK Homepage</h1>
                </div>
                <div style={{ textAlign: "center" }}>
                    <Button variant="outlined"
                        style={{ fontSize: "18px", color: "white", borderColor: "white" }}
                        component={Link} to="/login"
                    >Login</Button>
                    <span style={{ marginLeft: "30px" }}></span>
                    <Button variant="outlined"
                        style={{ fontSize: "18px", color: "white", borderColor: "white" }}
                        component={Link} to="/signup"
                    >Sign up</Button>
                </div>
            </div>
        </>
    )
}

export default Homepage;
