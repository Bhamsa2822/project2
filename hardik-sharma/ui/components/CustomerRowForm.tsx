import { Td, Button, Input, FormControl, InputLeftAddon, InputGroup, FormLabel } from "@chakra-ui/react";
import { Customer } from "./customer";
import { useState, ChangeEvent } from "react";

export interface CustomerRowFormProps {
    customer: Customer
    submit(customer: Customer): void
}

interface CustomerFormData {
    name: string
    address: string
    contactNo: string
}

export default function CustomerRowForm({ customer, submit }: CustomerRowFormProps) {
    const [formInputData, setFormInputData] = useState<CustomerFormData>({
        name: customer.customerDetails.name,
        address: customer.customerDetails.address,
        contactNo: customer.customerDetails.contactNo.toString()
    })

    const handleChange = (e: ChangeEvent<HTMLInputElement>): void => {
        setFormInputData({ ...formInputData, [e.target.name]: e.target.value })
    }

    const handleSubmit = (): void => {
        const cust: Customer = {
            id: customer.id,
            customerDetails: {
                name: formInputData.name,
                address: formInputData.address,
                contactNo: parseInt(formInputData.contactNo)
            }
        }
        submit(cust)
    }

    return (
        <>
            <Td>
                <FormControl isRequired>
                    <FormLabel>Name</FormLabel>
                    <Input
                        variant="filled"
                        color="chakra-body-bg._dark"
                        width="240px"
                        placeholder="Name"
                        type='text'
                        name="name"
                        defaultValue={formInputData.name}
                        onChange={handleChange}
                    />
                </FormControl>
            </Td>
            <Td>
                <FormControl isRequired>
                    <FormLabel>Address</FormLabel>
                    <Input
                        variant="filled"
                        color="chakra-body-bg._dark"
                        width="240px"
                        type='text'
                        placeholder="Address"
                        name="address"
                        defaultValue={formInputData.address}
                        onChange={handleChange}
                    />
                </FormControl>
            </Td>
            <Td>
                <FormControl isRequired>
                    <FormLabel>Contact No.</FormLabel>
                    <InputGroup>
                        <InputLeftAddon
                            children='+91'
                        />
                        <Input
                            variant="filled"
                            color="chakra-body-bg._dark"
                            width="240px"
                            type='number'
                            placeholder="Contact Number"
                            name="contactNo"
                            defaultValue={formInputData.contactNo}
                            onChange={handleChange}
                        />
                    </InputGroup>
                </FormControl>
            </Td>
            <Td>
                <Button
                    onClick={handleSubmit}
                    background="green.300">
                    {customer.id === "" ? "register" : "update"}
                </Button>
            </Td>
        </>
    )
}