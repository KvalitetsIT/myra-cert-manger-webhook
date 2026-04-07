package adaptors_test

/*
func TestMyraClientAdaptor_OnAdd(t *testing.T) {
	mockClient := new(client.MockMyraClient)
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)

	mockClient.On("OnAdd", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == record.ResolvedFQDN && r.Value == record.Key
	})).Return(dnsRecord, nil)

	result, err := adaptor.OnAdd(record)
	require.NoError(t, err, "expected no error from OnAdd")
	require.Equal(t, record.DNSName, result.DNSName, "DNSName should map back correctly")
	require.Equal(t, record.Key, result.Key, "Key should map back correctly")

	mockClient.AssertExpectations(t)
}

func TestMyraClientAdaptor_OnAdd_Error(t *testing.T) {
	mockClient := new(client.MockMyraClient)
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)
	clientErr := errors.New("client failed")

	mockClient.On("OnAdd", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == record.ResolvedFQDN && r.Value == record.Key
	})).Return(dnsRecord, clientErr)

	result, err := adaptor.OnAdd(record)
	require.Error(t, err, "expected error from OnAdd when client fails")
	require.Equal(t, clientErr, err, "error should match the client error")
	require.Equal(t, record.DNSName, result.DNSName, "DNSName should still map back on error")
	require.Equal(t, record.Key, result.Key, "Key should still map back on error")

	mockClient.AssertExpectations(t)
}

func TestMyraClientAdaptor_OnDelete(t *testing.T) {
	mockClient := new(client.MockMyraClient)
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)

	mockClient.On("OnDelete", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == record.ResolvedFQDN && r.Value == record.Key
	})).Return(dnsRecord, nil)

	result, err := adaptor.OnDelete(record)
	require.NoError(t, err, "expected no error from OnDelete")
	require.Equal(t, record.DNSName, result.DNSName, "DNSName should map back correctly")
	require.Equal(t, record.Key, result.Key, "Key should map back correctly")

	mockClient.AssertExpectations(t)
}

func TestMyraClientAdaptor_OnDelete_Error(t *testing.T) {
	mockClient := new(client.MockMyraClient)
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)
	clientErr := errors.New("delete failed")

	mockClient.On("OnDelete", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == record.ResolvedFQDN && r.Value == record.Key
	})).Return(dnsRecord, clientErr)

	result, err := adaptor.OnDelete(record)
	require.Error(t, err, "expected error from OnDelete when client fails")
	require.Equal(t, clientErr, err, "error should match the client error")
	require.Equal(t, record.DNSName, result.DNSName, "DNSName should still map back on error")
	require.Equal(t, record.Key, result.Key, "Key should still map back on error")

	mockClient.AssertExpectations(t)
}
*/