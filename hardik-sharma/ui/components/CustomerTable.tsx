import {
    Table,
    Thead,
    Tr,
    Th,
    TableContainer,
    Tbody,
    Button,
    Td,
} from '@chakra-ui/react'
import CustomerRow from './CustomerRow'
import { Customer } from './customer'
import { useState } from 'react'
import CustomerRowForm from './CustomerRowForm'

const defaultCustomer: Customer = {
    id: "",
    customerDetails: {
        name: "",
        address: "",
        contactNo: 0
    }
}

export interface CustomerTableProps {
    customers: Customer[]
    saveCustomer(customer: Customer): Promise<void>
    addCustomer(customer: Customer): Promise<void>
    onDelete(id: string): Promise<void>
}

export default function CustomerTable({ customers, saveCustomer, addCustomer, onDelete }: CustomerTableProps) {
    const [isAddForm, setIsAddForm] = useState(false)

    const showForm = () => {
        setIsAddForm(true)
    }

    const hideForm = () => {
        setIsAddForm(false)
    }

    const handleAdd = (c: Customer): Promise<void> => {
        return addCustomer(c).then(hideForm).catch((err) => { })
    }

    return (
        <TableContainer>
            <Table
                size="lg"
                variant="striped"
                colorScheme="green"
                backgroundColor=" #ffffbf"
                border="1px solid black"
            >
                <Thead>
                    <Tr>
                        <Th>ID</Th>
                        <Th>Name</Th>
                        <Th>Address</Th>
                        <Th>Contact Number</Th>
                        <Th>
                            <Button onClick={showForm}>register</Button>
                        </Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {isAddForm && (
                        <Tr>
                            <Td></Td>
                            <CustomerRowForm
                                customer={defaultCustomer}
                                submit={handleAdd}
                            />
                            <Td></Td>
                        </Tr>
                    )}
                    {customers.map((customer, index) =>
                        <CustomerRow
                            key={index}
                            customer={customer}
                            saveCustomer={saveCustomer}
                            onDelete={onDelete}
                        />
                    )}
                </Tbody>
            </Table>
        </TableContainer>
    )
}