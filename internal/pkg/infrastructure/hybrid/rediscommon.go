//
// SPDX-License-Identifier: Apache-2.0

package hybrid

import (
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	model "github.com/edgexfoundry/go-mod-core-contracts/v3/models"
)

// Add a new device profle
func (c *HybridClient) AddDeviceProfile(dp model.DeviceProfile) (model.DeviceProfile, errors.EdgeX) {
	return c.redisClient.AddDeviceProfile(dp)
}

// UpdateDeviceProfile updates a new device profile
func (c *HybridClient) UpdateDeviceProfile(dp model.DeviceProfile) errors.EdgeX {
	return c.redisClient.UpdateDeviceProfile(dp)
}

// DeviceProfileNameExists checks the device profile exists by name
func (c *HybridClient) DeviceProfileNameExists(name string) (bool, errors.EdgeX) {
	return c.redisClient.DeviceProfileNameExists(name)
}

// AddDeviceService adds a new device service
func (c *HybridClient) AddDeviceService(ds model.DeviceService) (model.DeviceService, errors.EdgeX) {
	return c.redisClient.AddDeviceService(ds)
}

// DeviceServiceByName gets a device service by name
func (c *HybridClient) DeviceServiceByName(name string) (deviceService model.DeviceService, edgeXerr errors.EdgeX) {
	return c.redisClient.DeviceServiceByName(name)
}

// DeviceServiceById gets a device service by id
func (c *HybridClient) DeviceServiceById(id string) (deviceService model.DeviceService, edgeXerr errors.EdgeX) {
	return c.redisClient.DeviceServiceById(id)
}

// DeleteDeviceServiceById deletes a device service by id
func (c *HybridClient) DeleteDeviceServiceById(id string) errors.EdgeX {
	return c.redisClient.DeleteDeviceServiceById(id)
}

// DeleteDeviceServiceByName deletes a device service by name
func (c *HybridClient) DeleteDeviceServiceByName(name string) errors.EdgeX {
	return c.redisClient.DeleteDeviceServiceByName(name)
}

// DeviceServiceNameExists checks the device service exists by name
func (c *HybridClient) DeviceServiceNameExists(name string) (bool, errors.EdgeX) {
	return c.redisClient.DeviceServiceNameExists(name)
}

// UpdateDeviceService updates a device service
func (c *HybridClient) UpdateDeviceService(ds model.DeviceService) errors.EdgeX {
	return c.redisClient.UpdateDeviceService(ds)
}

// DeviceProfileById gets a device profile by id
func (c *HybridClient) DeviceProfileById(id string) (deviceProfile model.DeviceProfile, err errors.EdgeX) {
	return c.redisClient.DeviceProfileById(id)
}

// DeviceProfileByName gets a device profile by name
func (c *HybridClient) DeviceProfileByName(name string) (deviceProfile model.DeviceProfile, edgeXerr errors.EdgeX) {
	return c.redisClient.DeviceProfileByName(name)
}

// DeleteDeviceProfileById deletes a device profile by id
func (c *HybridClient) DeleteDeviceProfileById(id string) errors.EdgeX {
	return c.redisClient.DeleteDeviceProfileById(id)
}

// DeleteDeviceProfileByName deletes a device profile by name
func (c *HybridClient) DeleteDeviceProfileByName(name string) errors.EdgeX {
	return c.redisClient.DeleteDeviceProfileByName(name)
}

// AllDeviceProfiles query device profiles with offset and limit
func (c *HybridClient) AllDeviceProfiles(offset int, limit int, labels []string) ([]model.DeviceProfile, errors.EdgeX) {
	return c.redisClient.AllDeviceProfiles(offset, limit, labels)
}

// DeviceProfilesByModel query device profiles with offset, limit and model
func (c *HybridClient) DeviceProfilesByModel(offset int, limit int, model string) ([]model.DeviceProfile, errors.EdgeX) {
	return c.redisClient.DeviceProfilesByModel(offset, limit, model)
}

// DeviceProfilesByManufacturer query device profiles with offset, limit and manufacturer
func (c *HybridClient) DeviceProfilesByManufacturer(offset int, limit int, manufacturer string) ([]model.DeviceProfile, errors.EdgeX) {
	return c.redisClient.DeviceProfilesByManufacturer(offset, limit, manufacturer)
}

// DeviceProfilesByManufacturerAndModel query device profiles with offset, limit, manufacturer and model
func (c *HybridClient) DeviceProfilesByManufacturerAndModel(offset int, limit int, manufacturer string, model string) ([]model.DeviceProfile, uint32, errors.EdgeX) {
	return c.redisClient.DeviceProfilesByManufacturerAndModel(offset, limit, manufacturer, model)
}

// AllDeviceServices returns multiple device services per query criteria, including
// offset: the number of items to skip before starting to collect the result set
// limit: The numbers of items to return
// labels: allows for querying a given object by associated user-defined labels
func (c *HybridClient) AllDeviceServices(offset int, limit int, labels []string) (deviceServices []model.DeviceService, edgeXerr errors.EdgeX) {
	return c.redisClient.AllDeviceServices(offset, limit, labels)
}

// Add a new device
func (c *HybridClient) AddDevice(d model.Device) (model.Device, errors.EdgeX) {
	return c.redisClient.AddDevice(d)
}

// DeleteDeviceById deletes a device by id
func (c *HybridClient) DeleteDeviceById(id string) errors.EdgeX {
	return c.redisClient.DeleteDeviceById(id)
}

// DeleteDeviceByName deletes a device by name
func (c *HybridClient) DeleteDeviceByName(name string) errors.EdgeX {
	return c.redisClient.DeleteDeviceByName(name)
}

// DevicesByServiceName query devices by offset, limit and name
func (c *HybridClient) DevicesByServiceName(offset int, limit int, name string) (devices []model.Device, edgeXerr errors.EdgeX) {
	return c.redisClient.DevicesByServiceName(offset, limit, name)
}

// DeviceIdExists checks the device existence by id
func (c *HybridClient) DeviceIdExists(id string) (bool, errors.EdgeX) {
	return c.redisClient.DeviceIdExists(id)
}

// DeviceNameExists checks the device existence by name
func (c *HybridClient) DeviceNameExists(name string) (bool, errors.EdgeX) {
	return c.redisClient.DeviceNameExists(name)
}

// DeviceById gets a device by id
func (c *HybridClient) DeviceById(id string) (device model.Device, edgeXerr errors.EdgeX) {
	return c.redisClient.DeviceById(id)
}

// DeviceByName gets a device by name
func (c *HybridClient) DeviceByName(name string) (device model.Device, edgeXerr errors.EdgeX) {
	return c.redisClient.DeviceByName(name)
}

// DevicesByProfileName query devices by offset, limit and profile name
func (c *HybridClient) DevicesByProfileName(offset int, limit int, profileName string) (devices []model.Device, edgeXerr errors.EdgeX) {
	return c.redisClient.DevicesByProfileName(offset, limit, profileName)
}

// Update a device
func (c *HybridClient) UpdateDevice(d model.Device) errors.EdgeX {
	return c.redisClient.UpdateDevice(d)
}

// AllDevices query the devices with offset, limit, and labels
func (c *HybridClient) AllDevices(offset int, limit int, labels []string) ([]model.Device, errors.EdgeX) {
	return c.redisClient.AllDevices(offset, limit, labels)
}

// AddProvisionWatcher adds a new provision watcher
func (c *HybridClient) AddProvisionWatcher(pw model.ProvisionWatcher) (model.ProvisionWatcher, errors.EdgeX) {
	return c.redisClient.AddProvisionWatcher(pw)
}

// ProvisionWatcherById gets a provision watcher by id
func (c *HybridClient) ProvisionWatcherById(id string) (provisionWatcher model.ProvisionWatcher, edgexErr errors.EdgeX) {
	return c.redisClient.ProvisionWatcherById(id)
}

// ProvisionWatcherByName gets a provision watcher by name
func (c *HybridClient) ProvisionWatcherByName(name string) (provisionWatcher model.ProvisionWatcher, edgexErr errors.EdgeX) {
	return c.redisClient.ProvisionWatcherByName(name)
}

// ProvisionWatchersByServiceName query provision watchers by offset, limit and service name
func (c *HybridClient) ProvisionWatchersByServiceName(offset int, limit int, name string) (provisionWatchers []model.ProvisionWatcher, edgexErr errors.EdgeX) {
	return c.redisClient.ProvisionWatchersByServiceName(offset, limit, name)
}

// ProvisionWatchersByProfileName query provision watchers by offset, limit and profile name
func (c *HybridClient) ProvisionWatchersByProfileName(offset int, limit int, name string) (provisionWatchers []model.ProvisionWatcher, edgexErr errors.EdgeX) {
	return c.redisClient.ProvisionWatchersByProfileName(offset, limit, name)
}

// AllProvisionWatchers query provision watchers with offset, limit and labels
func (c *HybridClient) AllProvisionWatchers(offset int, limit int, labels []string) (provisionWatchers []model.ProvisionWatcher, edgexErr errors.EdgeX) {
	return c.redisClient.AllProvisionWatchers(offset, limit, labels)
}

// DeleteProvisionWatcherByName deletes a provision watcher by name
func (c *HybridClient) DeleteProvisionWatcherByName(name string) errors.EdgeX {
	return c.redisClient.DeleteProvisionWatcherByName(name)
}

// Update a provision watcher
func (c *HybridClient) UpdateProvisionWatcher(pw model.ProvisionWatcher) errors.EdgeX {
	return c.redisClient.UpdateProvisionWatcher(pw)
}

// DeviceProfileCountByLabels returns the total count of Device Profiles with labels specified.  If no label is specified, the total count of all device profiles will be returned.
func (c *HybridClient) DeviceProfileCountByLabels(labels []string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceProfileCountByLabels(labels)
}

// DeviceProfileCountByManufacturer returns the count of Device Profiles associated with specified manufacturer
func (c *HybridClient) DeviceProfileCountByManufacturer(manufacturer string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceProfileCountByManufacturer(manufacturer)
}

// DeviceProfileCountByModel returns the count of Device Profiles associated with specified model
func (c *HybridClient) DeviceProfileCountByModel(model string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceProfileCountByModel(model)
}

// DeviceServiceCountByLabels returns the total count of Device Services with labels specified.  If no label is specified, the total count of all device services will be returned.
func (c *HybridClient) DeviceServiceCountByLabels(labels []string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceServiceCountByLabels(labels)
}

// DeviceCountByLabels returns the total count of Devices with labels specified.  If no label is specified, the total count of all devices will be returned.
func (c *HybridClient) DeviceCountByLabels(labels []string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceCountByLabels(labels)
}

// DeviceCountByProfileName returns the count of Devices associated with specified profile
func (c *HybridClient) DeviceCountByProfileName(profileName string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceCountByProfileName(profileName)
}

// DeviceCountByServiceName returns the count of Devices associated with specified service
func (c *HybridClient) DeviceCountByServiceName(serviceName string) (uint32, errors.EdgeX) {
	return c.redisClient.DeviceCountByServiceName(serviceName)
}

// ProvisionWatcherCountByLabels returns the total count of Provision Watchers with labels specified.  If no label is specified, the total count of all provision watchers will be returned.
func (c *HybridClient) ProvisionWatcherCountByLabels(labels []string) (uint32, errors.EdgeX) {
	return c.redisClient.ProvisionWatcherCountByLabels(labels)
}

// ProvisionWatcherCountByServiceName returns the count of Provision Watcher associated with specified service
func (c *HybridClient) ProvisionWatcherCountByServiceName(name string) (uint32, errors.EdgeX) {
	return c.redisClient.ProvisionWatcherCountByServiceName(name)
}

// ProvisionWatcherCountByProfileName returns the count of Provision Watcher associated with specified profile
func (c *HybridClient) ProvisionWatcherCountByProfileName(name string) (uint32, errors.EdgeX) {
	return c.redisClient.ProvisionWatcherCountByProfileName(name)
}

// AddInterval adds a new interval
func (c *HybridClient) AddInterval(interval model.Interval) (model.Interval, errors.EdgeX) {
	return c.redisClient.AddInterval(interval)
}

// IntervalByName gets a interval by name
func (c *HybridClient) IntervalByName(name string) (interval model.Interval, edgeXerr errors.EdgeX) {
	return c.redisClient.IntervalByName(name)
}

// IntervalById gets a interval by id
func (c *HybridClient) IntervalById(id string) (interval model.Interval, edgeXerr errors.EdgeX) {
	return c.redisClient.IntervalById(id)
}

// AllIntervals query intervals with offset and limit
func (c *HybridClient) AllIntervals(offset int, limit int) (intervals []model.Interval, edgeXerr errors.EdgeX) {
	return c.redisClient.AllIntervals(offset, limit)
}

// UpdateInterval updates a interval
func (c *HybridClient) UpdateInterval(interval model.Interval) errors.EdgeX {
	return c.redisClient.UpdateInterval(interval)
}

// DeleteIntervalByName deletes the interval by name
func (c *HybridClient) DeleteIntervalByName(name string) errors.EdgeX {
	return c.redisClient.DeleteIntervalByName(name)
}

// IntervalTotalCount returns the total count of Interval from the database
func (c *HybridClient) IntervalTotalCount() (uint32, errors.EdgeX) {
	return c.redisClient.IntervalTotalCount()
}

// IntervalActionTotalCount returns the total count of IntervalAction from the database
func (c *HybridClient) IntervalActionTotalCount() (uint32, errors.EdgeX) {
	return c.redisClient.IntervalActionTotalCount()
}

// AddIntervalAction adds a new intervalAction
func (c *HybridClient) AddIntervalAction(action model.IntervalAction) (model.IntervalAction, errors.EdgeX) {
	return c.redisClient.AddIntervalAction(action)
}

// AllIntervalActions query intervalActions with offset and limit
func (c *HybridClient) AllIntervalActions(offset int, limit int) (intervalActions []model.IntervalAction, edgeXerr errors.EdgeX) {
	return c.redisClient.AllIntervalActions(offset, limit)
}

// IntervalActionByName gets a intervalAction by name
func (c *HybridClient) IntervalActionByName(name string) (action model.IntervalAction, edgeXerr errors.EdgeX) {
	return c.redisClient.IntervalActionByName(name)
}

// IntervalActionsByIntervalName query intervalActions by offset, limit and intervalName
func (c *HybridClient) IntervalActionsByIntervalName(offset int, limit int, intervalName string) (actions []model.IntervalAction, edgeXerr errors.EdgeX) {
	return c.redisClient.IntervalActionsByIntervalName(offset, limit, intervalName)
}

// DeleteIntervalActionByName deletes the intervalAction by name
func (c *HybridClient) DeleteIntervalActionByName(name string) errors.EdgeX {
	return c.redisClient.DeleteIntervalActionByName(name)
}

// IntervalActionById gets a intervalAction by id
func (c *HybridClient) IntervalActionById(id string) (action model.IntervalAction, edgeXerr errors.EdgeX) {
	return c.redisClient.IntervalActionById(id)
}

// UpdateIntervalAction updates a intervalAction
func (c *HybridClient) UpdateIntervalAction(action model.IntervalAction) errors.EdgeX {
	return c.redisClient.UpdateIntervalAction(action)
}

// AddSubscription adds a new subscription
func (c *HybridClient) AddSubscription(subscription model.Subscription) (model.Subscription, errors.EdgeX) {
	return c.redisClient.AddSubscription(subscription)
}

// AllSubscriptions returns multiple subscriptions per query criteria, including
// offset: The number of items to skip before starting to collect the result set.
// limit: The maximum number of items to return.
func (c *HybridClient) AllSubscriptions(offset int, limit int) ([]model.Subscription, errors.EdgeX) {
	return c.redisClient.AllSubscriptions(offset, limit)
}

// SubscriptionsByCategory queries subscriptions by offset, limit and category
func (c *HybridClient) SubscriptionsByCategory(offset int, limit int, category string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	return c.redisClient.SubscriptionsByCategory(offset, limit, category)
}

// SubscriptionsByLabel queries subscriptions by offset, limit and label
func (c *HybridClient) SubscriptionsByLabel(offset int, limit int, label string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	return c.redisClient.SubscriptionsByLabel(offset, limit, label)
}

// SubscriptionsByReceiver queries subscriptions by offset, limit and receiver
func (c *HybridClient) SubscriptionsByReceiver(offset int, limit int, receiver string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	return c.redisClient.SubscriptionsByReceiver(offset, limit, receiver)
}

// SubscriptionById gets a subscription by id
func (c *HybridClient) SubscriptionById(id string) (subscription model.Subscription, edgexErr errors.EdgeX) {
	return c.redisClient.SubscriptionById(id)
}

// SubscriptionByName queries subscription by name
func (c *HybridClient) SubscriptionByName(name string) (subscription model.Subscription, edgeXerr errors.EdgeX) {
	return c.redisClient.SubscriptionByName(name)
}

// UpdateSubscription updates a new subscription
func (c *HybridClient) UpdateSubscription(subscription model.Subscription) errors.EdgeX {
	return c.redisClient.UpdateSubscription(subscription)
}

// DeleteSubscriptionByName deletes a subscription by name
func (c *HybridClient) DeleteSubscriptionByName(name string) errors.EdgeX {
	return c.redisClient.DeleteSubscriptionByName(name)
}

// SubscriptionsByCategoriesAndLabels queries subscriptions by offset, limit, categories and labels
func (c *HybridClient) SubscriptionsByCategoriesAndLabels(offset int, limit int, categories []string, labels []string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	return c.redisClient.SubscriptionsByCategoriesAndLabels(offset, limit, categories, labels)
}

// AddNotification adds a new notification
func (c *HybridClient) AddNotification(notification model.Notification) (model.Notification, errors.EdgeX) {
	return c.redisClient.AddNotification(notification)
}

// NotificationsByCategory queries notifications by offset, limit and category
func (c *HybridClient) NotificationsByCategory(offset int, limit int, category string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	return c.redisClient.NotificationsByCategory(offset, limit, category)
}

// NotificationsByLabel queries notifications by offset, limit and label
func (c *HybridClient) NotificationsByLabel(offset int, limit int, label string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	return c.redisClient.NotificationsByLabel(offset, limit, label)
}

// NotificationById gets a notification by id
func (c *HybridClient) NotificationById(id string) (notification model.Notification, edgexErr errors.EdgeX) {
	return c.redisClient.NotificationById(id)
}

// NotificationsByStatus queries notifications by offset, limit and status
func (c *HybridClient) NotificationsByStatus(offset int, limit int, status string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	return c.redisClient.NotificationsByStatus(offset, limit, status)
}

// NotificationsByTimeRange query notifications by time range, offset, and limit
func (c *HybridClient) NotificationsByTimeRange(start int, end int, offset int, limit int) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	return c.redisClient.NotificationsByTimeRange(start, end, offset, limit)
}

// NotificationsByCategoriesAndLabels queries notifications by offset, limit, categories and labels
func (c *HybridClient) NotificationsByCategoriesAndLabels(offset int, limit int, categories []string, labels []string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	return c.redisClient.NotificationsByCategoriesAndLabels(offset, limit, categories, labels)
}

// NotificationCountByCategory returns the count of Notification associated with specified category from the database
func (c *HybridClient) NotificationCountByCategory(category string) (uint32, errors.EdgeX) {
	return c.redisClient.NotificationCountByCategory(category)
}

// NotificationCountByLabel returns the count of Notification associated with specified label from the database
func (c *HybridClient) NotificationCountByLabel(label string) (uint32, errors.EdgeX) {
	return c.redisClient.NotificationCountByLabel(label)
}

// NotificationCountByStatus returns the count of Notification associated with specified status from the database
func (c *HybridClient) NotificationCountByStatus(status string) (uint32, errors.EdgeX) {
	return c.redisClient.NotificationCountByStatus(status)
}

// NotificationCountByTimeRange returns the count of Notification from the database within specified time range
func (c *HybridClient) NotificationCountByTimeRange(start int, end int) (uint32, errors.EdgeX) {
	return c.redisClient.NotificationCountByTimeRange(start, end)
}

// NotificationCountByCategoriesAndLabels returns the count of Notification associated with specified categories and labels from the database
func (c *HybridClient) NotificationCountByCategoriesAndLabels(categories []string, labels []string) (uint32, errors.EdgeX) {
	return c.redisClient.NotificationCountByCategoriesAndLabels(categories, labels)
}

// SubscriptionTotalCount returns the total count of Subscription from the database
func (c *HybridClient) SubscriptionTotalCount() (uint32, errors.EdgeX) {
	return c.redisClient.SubscriptionTotalCount()
}

// SubscriptionCountByCategory returns the count of Subscription associated with specified category from the database
func (c *HybridClient) SubscriptionCountByCategory(category string) (uint32, errors.EdgeX) {
	return c.redisClient.SubscriptionCountByCategory(category)
}

// SubscriptionCountByLabel returns the count of Subscription associated with specified label from the database
func (c *HybridClient) SubscriptionCountByLabel(label string) (uint32, errors.EdgeX) {
	return c.redisClient.SubscriptionCountByLabel(label)
}

// SubscriptionCountByReceiver returns the count of Subscription associated with specified receiver from the database
func (c *HybridClient) SubscriptionCountByReceiver(receiver string) (uint32, errors.EdgeX) {
	return c.redisClient.SubscriptionCountByReceiver(receiver)
}

// TransmissionTotalCount returns the total count of Transmission from the database
func (c *HybridClient) TransmissionTotalCount() (uint32, errors.EdgeX) {
	return c.redisClient.TransmissionTotalCount()
}

// TransmissionCountBySubscriptionName returns the count of Transmission associated with specified subscription name from the database
func (c *HybridClient) TransmissionCountBySubscriptionName(subscriptionName string) (uint32, errors.EdgeX) {
	return c.redisClient.TransmissionCountBySubscriptionName(subscriptionName)
}

// TransmissionCountByStatus returns the count of Transmission associated with specified status name from the database
func (c *HybridClient) TransmissionCountByStatus(status string) (uint32, errors.EdgeX) {
	return c.redisClient.TransmissionCountByStatus(status)
}

// TransmissionCountByTimeRange returns the count of Transmission from the database within specified time range
func (c *HybridClient) TransmissionCountByTimeRange(start int, end int) (uint32, errors.EdgeX) {
	return c.redisClient.TransmissionCountByTimeRange(start, end)
}

// DeleteNotificationById deletes a notification by id
func (c *HybridClient) DeleteNotificationById(id string) errors.EdgeX {
	return c.redisClient.DeleteNotificationById(id)
}

// UpdateNotification updates a notification
func (c *HybridClient) UpdateNotification(n model.Notification) errors.EdgeX {
	return c.redisClient.UpdateNotification(n)
}

// AddTransmission adds a new transmission
func (c *HybridClient) AddTransmission(t model.Transmission) (model.Transmission, errors.EdgeX) {
	return c.redisClient.AddTransmission(t)
}

// UpdateTransmission updates a transmission
func (c *HybridClient) UpdateTransmission(trans model.Transmission) errors.EdgeX {
	return c.redisClient.UpdateTransmission(trans)
}

// TransmissionById gets a transmission by id
func (c *HybridClient) TransmissionById(id string) (trans model.Transmission, edgexErr errors.EdgeX) {
	return c.redisClient.TransmissionById(id)
}

// TransmissionsByTimeRange query transmissions by time range, offset, and limit
func (c *HybridClient) TransmissionsByTimeRange(start int, end int, offset int, limit int) (transmissions []model.Transmission, err errors.EdgeX) {
	return c.redisClient.TransmissionsByTimeRange(start, end, offset, limit)
}

// AllTransmissions returns multiple transmissions per query criteria, including
// offset: The number of items to skip before starting to collect the result set.
// limit: The maximum number of items to return.
func (c *HybridClient) AllTransmissions(offset int, limit int) ([]model.Transmission, errors.EdgeX) {
	return c.redisClient.AllTransmissions(offset, limit)
}

// TransmissionsByStatus queries transmissions by offset, limit and status
func (c *HybridClient) TransmissionsByStatus(offset int, limit int, status string) (transmissions []model.Transmission, err errors.EdgeX) {
	return c.redisClient.TransmissionsByStatus(offset, limit, status)
}

// TransmissionsBySubscriptionName queries transmissions by offset, limit and subscription name
func (c *HybridClient) TransmissionsBySubscriptionName(offset int, limit int, subscriptionName string) (transmissions []model.Transmission, err errors.EdgeX) {
	return c.redisClient.TransmissionsBySubscriptionName(offset, limit, subscriptionName)
}

// TransmissionsByNotificationId queries transmissions by offset, limit and notification id
func (c *HybridClient) TransmissionsByNotificationId(offset int, limit int, id string) (transmissions []model.Transmission, err errors.EdgeX) {
	return c.redisClient.TransmissionsByNotificationId(offset, limit, id)
}

// TransmissionCountByNotificationId returns the count of Transmission associated with specified notification id from the database
func (c *HybridClient) TransmissionCountByNotificationId(id string) (uint32, errors.EdgeX) {
	return c.redisClient.TransmissionCountByNotificationId(id)
}
