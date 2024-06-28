import {
    Tr,
    Td,
    Button,
} from '@chakra-ui/react'
import { Customer } from './customer'
import { useState } from 'react'
import CustomerRowForm from './CustomerRowForm'

export interface CustomerRowProps {
    customer: Customer
    saveCustomer(customer: Customer): Promise<void>
    onDelete(id: string): Promise<void>
}

export default function CustomerRow({ customer, saveCustomer, onDelete }: CustomerRowProps) {
    const [isFormRow, setIsFormRow] = useState(false)

    const showForm = (): void => {
        setIsFormRow(true)
    }

    const showData = (): void => {
        setIsFormRow(false)
    }


    function handleSave(customer: Customer): Promise<void> {
        return saveCustomer(customer)
            .then(showData).
            catch((err) => { })
    }

    if (isFormRow) {
        return <FormRow customer={customer} saveCustomer={handleSave} />
    }

    return <DataRow customer={customer} setShowFormTrue={showForm} onDelete={onDelete} />
}

interface RowDataProps {
    customer: Customer
    setShowFormTrue(): void
    onDelete(id: string): void
}

function DataRow({ customer, setShowFormTrue, onDelete }: RowDataProps) {
    const handleDelete = () => {
        onDelete(customer.id)
    }

    return (
        <Tr>
            <Td>{customer.id}</Td>
            <Td>{customer.customerDetails.name}</Td>
            <Td>{customer.customerDetails.address}</Td>
            <Td>{customer.customerDetails.contactNo}</Td>
            <Td>
                <Button backgroundColor="yellow.400" onClick={setShowFormTrue}>Edit</Button>
                <Button backgroundColor="red.400" onClick={handleDelete} >delete</Button>
            </Td>
        </Tr>
    )
}

interface FormRowProps {
    customer: Customer
    saveCustomer(customer: Customer): Promise<void>
}

export function FormRow({ customer, saveCustomer }: FormRowProps) {
    return (
        <Tr>
            <Td>{customer.id}</Td>
            <CustomerRowForm
                customer={customer}
                submit={saveCustomer}
            />
        </Tr>
    )
}