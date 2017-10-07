class OpSpecNodeApiClient {
  async startOp(startOpReq) {
    // @TODO: don't assume local node
    return fetch('http://localhost:42224/ops/starts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(startOpReq),
    })
      .then(response => (response.text()))
  }
}

// export as singleton
export default new OpSpecNodeApiClient();
