import apiClient from './apiClient.js'

const PAGE_LIMIT = 100

function validateListResponse(response, resource, expectedPage, expectedTotal) {
  const body = response.data
  const pagination = body?.pagination
  if (
    body?.success !== true ||
    !Array.isArray(body.data) ||
    !Number.isSafeInteger(pagination?.page) ||
    !Number.isSafeInteger(pagination?.limit) ||
    !Number.isSafeInteger(pagination?.total) ||
    pagination.page !== expectedPage ||
    pagination.limit !== PAGE_LIMIT ||
    pagination.total < 0 ||
    (expectedTotal !== undefined && pagination.total !== expectedTotal)
  ) {
    throw new Error(`The ${resource} API returned an unexpected paginated response.`)
  }
  return { records: body.data, total: pagination.total }
}

async function fetchAllPages(resource, signal) {
  const records = []
  const seenIDs = new Set()
  let total
  let page = 1

  while (total === undefined || records.length < total) {
    const response = await apiClient.get(`/${resource}`, {
      params: { page, limit: PAGE_LIMIT },
      signal,
    })
    const result = validateListResponse(response, resource, page, total)
    total = result.total

    const expectedLength = Math.min(PAGE_LIMIT, total - records.length)
    if (result.records.length !== expectedLength) {
      throw new Error(
        `The ${resource} records changed while loading. Refresh the dashboard and try again.`,
      )
    }

    for (const record of result.records) {
      const id = record?.[resource === 'customers' ? 'customer_id' : 'bill_id']
      if (!Number.isSafeInteger(id) || seenIDs.has(id)) {
        throw new Error(
          `The ${resource} API returned incomplete or duplicate records. Refresh the dashboard and try again.`,
        )
      }
      seenIDs.add(id)
    }
    records.push(...result.records)
    page += 1
  }

  if (records.length !== total || seenIDs.size !== total) {
    throw new Error(
      `Could not verify the complete ${resource} list. Refresh the dashboard and try again.`,
    )
  }

  return records
}

function validateBillAmounts(bills) {
  for (const bill of bills) {
    for (const field of ['amount', 'amount_received', 'amount_pending']) {
      if (!Number.isSafeInteger(bill[field])) {
        throw new Error(`A bill is missing a valid ${field} value.`)
      }
    }
    if (!Number.isSafeInteger(bill.bill_id) || !Number.isSafeInteger(bill.customer_id)) {
      throw new Error('A bill is missing a valid bill or customer ID.')
    }
  }
}

export async function getDashboardRecords(signal) {
  const [customers, bills] = await Promise.all([
    fetchAllPages('customers', signal),
    fetchAllPages('bills', signal),
  ])
  validateBillAmounts(bills)
  return { customers, bills }
}
