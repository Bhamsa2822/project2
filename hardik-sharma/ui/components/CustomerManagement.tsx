import { useEffect, useState } from 'react'
import { deleteCustomer, createCustomer, fetchAllCustomers, updateCustomer } from '@/api'
import { Customer } from './customer'
import { Alert, AlertDescription, AlertIcon, AlertTitle, Box, Center } from '@chakra-ui/react'
import CustomerTable from './CustomerTable'

export default function CustomerManagement() {
    const [customers, setCustomers] = useState<Customer[]>([])
    const [errMsgFetch, setErrMsgFetch] = useState<String>('')
    const [successMsg, setSuccessMsg] = useState<string>("")
    const [errMsg, setErrMsg] = useState<string>("")

    function loadCustomerList() {
        fetchAllCustomers()
            .then((customers) => { setCustomers(customers) })
            .catch((err) => {
                setErrMsgFetch(err.message)
            })
    }

    useEffect(() => {
        loadCustomerList()

        const socket = new WebSocket("ws://127.0.0.1:8080/ws")
        console.log("attempting to connect ......")

        socket.onopen = () => {
            console.log("websocket connected")
        }

        socket.onmessage = (msg) => {
            const message = msg.data
            const customerList = JSON.parse(message)
            setCustomers(customerList)
        }

        socket.onclose = () => {
            console.log("connection closed")
        }

        socket.onerror = (err) => {
            console.log(err)
        }

        return () => {
            socket.close()
        }
    }, [])

    function addCustomer(cust: Customer): Promise<void> {
        const id = Math.random().toString(36).substring(2, 4)
        cust.id = id

        return createCustomer(cust)
            .then((res) => {
                setErrMsg("")
                setSuccessMsg(res)
                setTimeout(() => {
                    setSuccessMsg("")
                }, 1000)
            })
            .catch((err) => {
                setErrMsg(err.response.data)
                setSuccessMsg("")
                throw err
            })
    }

    function saveCustomer(cust: Customer): Promise<void> {
        return updateCustomer(cust).then((res) => {
            setErrMsg("")
            setSuccessMsg(res)
            setTimeout(() => {
                setSuccessMsg("")
            }, 1000)
        })
            .catch((err) => {
                setErrMsg(err.response.data)
                setSuccessMsg("")
                throw err
            })
    }

    const onDelete = (id: string): Promise<void> => {
        return deleteCustomer(id)
            .then(() => {
                setErrMsg("")
            })
            .catch((err) => {
                setErrMsg(err.response.data)
            })

    }

    return (
        <>
            <Center>
                <Box width="250px">
                    {errMsg && (
                        <Alert status='error'>
                            <AlertIcon />
                            <AlertDescription>{errMsg}</AlertDescription>
                        </Alert>
                    )}
                    {successMsg && (
                        <Alert status='success'>
                            <AlertIcon />
                            <AlertDescription>{successMsg}</AlertDescription>
                        </Alert>
                    )}
                    {errMsgFetch && (
                        <Alert status='error'>
                            <AlertIcon />
                            <AlertTitle>Something went wrong!</AlertTitle>
                            <AlertDescription>{errMsgFetch}</AlertDescription>
                        </Alert>
                    )}
                </Box>
            </Center>
            <CustomerTable
                customers={customers}
                saveCustomer={saveCustomer}
                addCustomer={addCustomer}
                onDelete={onDelete}
            />
        </>
    )
}
