import React from 'react';
import Input from '@mui/material/Input';
import InputLabel from '@mui/material/InputLabel';
import InputAdornment from '@mui/material/InputAdornment';
import AccountCircle from '@mui/icons-material/AccountCircle';
import LockIcon from '@mui/icons-material/Lock';
import { Button } from '@mui/material';
import { Link } from 'react-router-dom';
import { Formik } from 'formik';

function Login() {
    return (
        <>
            <div className='inputStyle'>
                <h1>COVID BOOK</h1>

                <Formik
                    initialValues={{ email: '', password: '' }}
                    onSubmit={(values) => {
                        console.log(values);
                    }}
                >
                    {({
                        values,
                        handleChange,
                        handleSubmit,
                    }) => (
                        <form onSubmit={handleSubmit}>
                            <InputLabel htmlFor="input-with-icon-adornment">
                            </InputLabel>
                            <Input
                                onChange={handleChange}
                                value={values.email}
                                inputProps={{ style: { fontSize: 40, fontFamily: 'Dongle', color: 'white', width: '300px' } }}
                                name='email'
                                placeholder='Email'
                                type='email'
                                startAdornment={
                                    <InputAdornment position="start">
                                        <AccountCircle style={{ fontSize: '32px', color: 'white' }} />
                                    </InputAdornment>
                                }
                            />
                            {/* make space between email and password input */}
                            <div style={{ marginTop: '20px' }}></div>

                            {/* input password */}
                            <InputLabel htmlFor="input-with-icon-adornment">

                            </InputLabel>
                            <Input
                                onChange={handleChange}
                                value={values.password}
                                inputProps={{ style: { fontSize: 40, fontFamily: 'Dongle', color: 'white', width: '300px' } }}
                                type='password'
                                name='password'
                                placeholder='Password'
                                startAdornment={
                                    <InputAdornment position="start">
                                        <LockIcon style={{ fontSize: '32px', color: 'white' }} />
                                    </InputAdornment>
                                }
                            />

                            {/* make space between email and password input */}
                            <div style={{ marginTop: '40px' }}></div>

                            {/* submit button */}
                            <Button variant="outlined"
                                style={{ fontSize: "18px", color: "white", borderColor: "white" }}
                                type="submit"
                            >Login</Button>

                            <span style={{ marginLeft: "30px" }}></span>
                            <Button variant="outlined"
                                style={{ fontSize: "18px", color: "white", borderColor: "white" }}
                                component={Link} to="/signup"
                            >Sign up</Button>
                            <div style={{ marginTop: "20px" }}>
                                <Button variant="outlined"
                                    style={{ fontSize: "18px", color: "white", borderColor: "white" }}
                                    component={Link} to="/records"
                                >Records</Button>
                            </div>
                            {/* console log values */}
                            <pre>{console.log(values)}</pre>
                        </form>
                    )}
                </Formik>
            </div>
        </>
    )
}
export default Login;
