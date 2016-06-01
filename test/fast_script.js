const startDate = Date.now()

const main = function() {
  const endDate = Date.now()

  process.stdout.write(JSON.stringify({
    StartDate: startDate,
    EndDate: endDate,
  }))
}

main()
